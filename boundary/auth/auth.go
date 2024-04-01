package auth

import (
	"context"
	"time"

	"github.com/go-pkgz/auth"
	"github.com/go-pkgz/auth/avatar"
	"github.com/go-pkgz/auth/provider"
	"github.com/go-pkgz/auth/token"
	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
	"gitlab.com/sfz.aalen/hackwerk/dotinder/crypto"
	"gitlab.com/sfz.aalen/hackwerk/dotinder/entity"
)

func NewAuthOptions(repo *entity.Repository) *auth.Opts {
	// define options
	options := auth.Opts{
		SecretReader: token.SecretFunc(func(id string) (string, error) { // secret key for JWT
			return "secret", nil
		}),
		DisableXSRF:    true,            // needed for development on localhost
		TokenDuration:  time.Minute * 5, // token expires in 5 minutes
		CookieDuration: time.Hour * 24,  // cookie expires in 1 day and will enforce re-login
		Issuer:         "dotinder",
		URL:            "http://localhost:8080",
		AvatarStore:    avatar.NewLocalFS("/tmp"),
		BasicAuthChecker: func(user, passwd string) (bool, token.User, error) {
			tx := repo.Pool.MustBegin()
			dbUser, err := repo.FindPasswordUser(tx, user)
			if err != nil {
				if err = tx.Rollback(); err != nil {
					return false, token.User{}, err
				}
				return false, token.User{}, err
			}
			if err = tx.Rollback(); err != nil {
				return false, token.User{}, err
			}
			ok, err := crypto.ComparePasswordAndHash(passwd, dbUser.Password)
			if err != nil {
				return false, token.User{}, err
			}
			return ok, token.User{ID: dbUser.Uuid.String(), Name: dbUser.Username}, nil
		},
	}

	return &options
}

func NewAuthService(options *auth.Opts, repo *entity.Repository, ctx context.Context) *auth.Service {
	// create auth service with providers
	service := auth.NewService(*options)
	msgTemplate := "http://localhost:8080/auth/matrix/login?token={{.Token}}"
	service.AddVerifProvider("matrix", msgTemplate, provider.SenderFunc(func(address string, text string) error {
		log.Ctx(ctx).Info().Msgf("sending message to %s: %s", address, text)
		return nil
	}))

	service.AddDirectProvider("local", provider.CredCheckerFunc(func(user, password string) (ok bool, err error) {
		tx := repo.Pool.MustBegin()
		dbUser, err := repo.FindPasswordUser(tx, user)
		if err != nil {
			if err = tx.Rollback(); err != nil {
				return false, err
			}
			return false, err
		}
		if err = tx.Rollback(); err != nil {
			return false, err
		}
		ok, err = crypto.ComparePasswordAndHash(password, dbUser.Password)
		return ok, err
	}))

	return service
}

func NewAuthRouter(service *auth.Service, router *mux.Router) *mux.Router {
	// setup auth routes
	authRoutes, avaRoutes := service.Handlers()
	router.PathPrefix("/auth").Handler(authRoutes)  // add auth handlers
	router.PathPrefix("/avatar").Handler(avaRoutes) // add avatar handlers

	// retrieve auth middleware
	m := service.Middleware()
	router.Use(m.Trace)

	authRouter := router.NewRoute().Subrouter()
	authRouter.Use(m.Auth)
	authRouter.Use(m.Trace)

	return authRouter
}
