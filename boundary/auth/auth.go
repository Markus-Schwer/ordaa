package auth

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
	"gitlab.com/sfz.aalen/hackwerk/dotinder/crypto"
	"gitlab.com/sfz.aalen/hackwerk/dotinder/entity"
)

// Create the JWT key used to create the signature
var jwtKey = []byte("my_secret_key")

// Create a struct to read the username and password from the request body
type Credentials struct {
	Password string `json:"password"`
	Username string `json:"username"`
}

// Create a struct that will be encoded to a JWT.
// We add jwt.RegisteredClaims as an embedded type, to provide fields like expiry time
type Claims struct {
	Username string `json:"username"`
	UserUuid string `json:"user_uuid"`
	jwt.RegisteredClaims
}

type AuthService struct {
	ctx  context.Context
	repo *entity.Repository
}

func NewAuthService(ctx context.Context, repo *entity.Repository) *AuthService {
	return &AuthService{ctx: ctx, repo: repo}
}

func (a *AuthService) Signin(tx *sqlx.Tx, creds *Credentials) (*jwt.Token, error) {
	dbUser, err := a.repo.FindPasswordUser(tx, creds.Username)
	if err != nil {
		return nil, err
	}
	ok, err := crypto.ComparePasswordAndHash(creds.Password, dbUser.Password)
	if err != nil || !ok {
		return nil, errors.New("invalid credentials")
	}

	// Declare the expiration time of the token
	// here, we have kept it as 5 minutes
	expirationTime := time.Now().Add(5 * time.Minute)
	// Create the JWT claims, which includes the username and expiry time
	claims := &Claims{
		Username: creds.Username,
		UserUuid: dbUser.Uuid.String(),
		RegisteredClaims: jwt.RegisteredClaims{
			Subject: creds.Username,
			// In JWT, the expiry time is expressed as unix milliseconds
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	// Declare the token with the algorithm used for signing, and the claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	log.Ctx(a.ctx).Info().Msgf("token %v", token)

	// Finally, we set the client cookie for "token" as the JWT we just generated
	// we also set an expiry time which is the same as the token itself
	return token, nil
}

func SignToken(token *jwt.Token) (string, error) {
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", fmt.Errorf("error while signing token: %w", err)
	}
	return tokenString, nil
}

func SetJwtCookie(token *jwt.Token, w http.ResponseWriter, r *http.Request) error {
	// Create the JWT string
	tokenString, err := SignToken(token)
	if err != nil {
		return err
	}

	expoirationTime, err := token.Claims.GetExpirationTime()
	if err != nil {
		return fmt.Errorf("error while getting expiration time: %w", err)
	}

	http.SetCookie(w, &http.Cookie{
		Name:    "token",
		Value:   tokenString,
		Expires: expoirationTime.Time,
	})

	return nil
}

func (a *AuthService) CheckAuthCookie(r *http.Request) (*jwt.Token, bool) {
	// We can obtain the session token from the requests cookies, which come with every request
	c, err := r.Cookie("token")
	if err != nil {
		if err == http.ErrNoCookie {
			// If the cookie is not set, return an unauthorized status
			return nil, false
		}
		// For any other type of error, return a bad request status
		return nil, false
	}

	// Get the JWT string from the cookie
	tknStr := c.Value

	// Initialize a new instance of `Claims`
	claims := &Claims{}

	// Parse the JWT string and store the result in `claims`.
	// Note that we are passing the key in this method as well. This method will return an error
	// if the token is invalid (if it has expired according to the expiry time we set on sign in),
	// or if the signature does not match
	tkn, err := jwt.ParseWithClaims(tknStr, claims, func(token *jwt.Token) (any, error) {
		return jwtKey, nil
	})
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			return nil, false
		}
		return nil, false
	}
	if !tkn.Valid {
		return nil, false
	}

	return tkn, true
}

func (a *AuthService) Refresh(w http.ResponseWriter, r *http.Request) error {
	token, ok := a.CheckAuthCookie(r)
	if !ok {
		return errors.New("not authenticated")
	}

	expiresAt, err := token.Claims.GetExpirationTime()
	if err != nil {
		return err
	}

	// We ensure that a new token is not issued until enough time has elapsed
	// In this case, a new token will only be issued if the old token is within
	// 30 seconds of expiry. Otherwise, return a bad request status
	if time.Until(expiresAt.Time) > 30*time.Second {
		return err
	}

	// Now, create a new token for the current use, with a renewed expiration time
	expirationTime := time.Now().Add(5 * time.Minute)
	// Initialize a new instance of `Claims`
	claims := &Claims{}
	claims.ExpiresAt = jwt.NewNumericDate(expirationTime)
	newToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	err = SetJwtCookie(newToken, w, r)
	if err != nil {
		return err
	}
	return nil
}

func (a *AuthService) Logout(w http.ResponseWriter, r *http.Request) {
	// immediately clear the token cookie
	http.SetCookie(w, &http.Cookie{
		Name:    "token",
		Expires: time.Now(),
	})
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func NewAuthRouter(auth *AuthService, router *mux.Router) *mux.Router {
	authRouter := router.NewRoute().Subrouter()
	authRouter.Use(mux.MiddlewareFunc(func(handler http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_, ok := auth.CheckAuthCookie(r)
			if !ok {
				auth.Logout(w, r)
				http.Redirect(w, r, "/login", http.StatusSeeOther)
			}
			auth.Refresh(w, r)
			handler.ServeHTTP(w, r)
		})
	}))

	return authRouter
}
