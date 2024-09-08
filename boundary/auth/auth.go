package auth

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"gitlab.com/sfz.aalen/hackwerk/dotinder/crypto"
	"gitlab.com/sfz.aalen/hackwerk/dotinder/entity"
	"gorm.io/gorm"
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
	repo entity.Repository
}

func NewAuthService(ctx context.Context, repo entity.Repository) *AuthService {
	return &AuthService{ctx: ctx, repo: repo}
}

func (a *AuthService) Signin(tx *gorm.DB, creds *Credentials) (*jwt.Token, error) {
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
		UserUuid: dbUser.UserUuid.String(),
		RegisteredClaims: jwt.RegisteredClaims{
			// In JWT, the expiry time is expressed as unix milliseconds
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	// Declare the token with the algorithm used for signing, and the claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

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

func SetJwtCookie(token *jwt.Token, c echo.Context) error {
	// Create the JWT string
	tokenString, err := SignToken(token)
	if err != nil {
		return err
	}

	expoirationTime, err := token.Claims.GetExpirationTime()
	if err != nil {
		return fmt.Errorf("error while getting expiration time: %w", err)
	}

	c.SetCookie(&http.Cookie{
		Name:    "token",
		Value:   tokenString,
		Expires: expoirationTime.Time,
	})

	return nil
}

func (a *AuthService) CheckTokenString(tknStr string) (*jwt.Token, bool) {
	if tknStr == "" {
		return nil, false
	}

	claims := &Claims{}

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

func (a *AuthService) CheckAuthHeader(c echo.Context) (*jwt.Token, bool) {
	authHeader := c.Request().Header.Get("Authorization")
	if authHeader == "" {
		return nil, false
	}
	splitAuthHeader := strings.Split(authHeader, " ")
	if len(splitAuthHeader) != 2 {
		return nil, false
	}

	tknStr := splitAuthHeader[1]
	log.Ctx(a.ctx).Info().Msgf("token string %s", tknStr)
	return a.CheckTokenString(tknStr)
}

func (a *AuthService) CheckAuthCookie(c echo.Context) (*jwt.Token, bool) {
	cookie, err := c.Cookie("token")
	if err != nil {
		if err == http.ErrNoCookie {
			return nil, false
		}
		return nil, false
	}

	tknStr := cookie.Value
	return a.CheckTokenString(tknStr)
}

func (a *AuthService) Refresh(c echo.Context) error {
	token, ok := a.CheckAuthCookie(c)
	if !ok {
		return errors.New("not authenticated")
	}

	oldClaims, ok := token.Claims.(*Claims)
	if !ok {
		return errors.New("could not cast claims")
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
	claims := &Claims{Username: oldClaims.Username, UserUuid: oldClaims.UserUuid}
	claims.ExpiresAt = jwt.NewNumericDate(expirationTime)
	newToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	err = SetJwtCookie(newToken, c)
	if err != nil {
		return err
	}
	return nil
}

func (a *AuthService) Logout(c echo.Context) {
	// immediately clear the token cookie
	c.SetCookie(&http.Cookie{
		Name:    "token",
		Expires: time.Now(),
	})
}

func AuthMiddleware(auth *AuthService, unauthorizedHandler func(echo.Context, error) error) echo.MiddlewareFunc {
	return echojwt.WithConfig(echojwt.Config{
		SigningKey:  jwtKey,
		ErrorHandler: unauthorizedHandler,
		NewClaimsFunc: func(c echo.Context) jwt.Claims {
			return new(Claims)
		},
		//SigningMethod: jwt.SigningMethodRS256.Name,
	})
}
