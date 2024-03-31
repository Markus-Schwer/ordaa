package rest

import (
	"context"
	"net/http"
	"time"

	"github.com/go-pkgz/auth"
	"github.com/go-pkgz/auth/avatar"
	"github.com/go-pkgz/auth/provider"
	"github.com/go-pkgz/auth/token"
	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
	"gitlab.com/sfz.aalen/hackwerk/dotinder/crypto"
	"gitlab.com/sfz.aalen/hackwerk/dotinder/entity"

	_ "github.com/lib/pq"
)

const AddressKey = "ADDRESS"

type RestBoundary struct {
	ctx  context.Context
	repo *entity.Repository
}

func NewRestBoundary(ctx context.Context, repo *entity.Repository) *RestBoundary {
	return &RestBoundary{ctx: ctx, repo: repo}
}

func (server *RestBoundary) Start() {
	// define options
	options := auth.Opts{
		SecretReader: token.SecretFunc(func(id string) (string, error) { // secret key for JWT
			return "secret", nil
		}),
		TokenDuration:  time.Minute * 5, // token expires in 5 minutes
		CookieDuration: time.Hour * 24,  // cookie expires in 1 day and will enforce re-login
		Issuer:         "dotinder",
		URL:            "http://localhost:8080",
		AvatarStore:    avatar.NewLocalFS("/tmp"),
		BasicAuthChecker: func(user, passwd string) (bool, token.User, error) {
			tx := server.repo.Pool.MustBegin()
			dbUser, err := server.repo.FindPasswordUser(tx, user)
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

	// create auth service with providers
	service := auth.NewService(options)
	service.AddDirectProviderWithUserIDFunc("local", provider.CredCheckerFunc(func(user, password string) (ok bool, err error) {
		tx := server.repo.Pool.MustBegin()
		dbUser, err := server.repo.FindPasswordUser(tx, user)
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
	}), provider.UserIDFunc(func(user string, r *http.Request) string {
		tx := server.repo.Pool.MustBegin()
		dbUser, err := server.repo.FindPasswordUser(tx, user)
		if err != nil {
			if err = tx.Rollback(); err != nil {
				log.Ctx(server.ctx).Error().Err(err).Msg(err.Error())
				return user
			}
			log.Ctx(server.ctx).Error().Err(err).Msg(err.Error())
			return user
		}
		if err = tx.Rollback(); err != nil {
			log.Ctx(server.ctx).Error().Err(err).Msg(err.Error())
			return user
		}
		return dbUser.UserUuid.String()
	}))

	// retrieve auth middleware
	m := service.Middleware()

	router := mux.NewRouter()
	authRouter := router.NewRoute().Subrouter()
	authRouter.Use(m.Trace)
	authRouter.Use(m.Auth)

	// setup auth routes
	authRoutes, avaRoutes := service.Handlers()
	router.PathPrefix("/auth").Handler(authRoutes)  // add auth handlers
	router.PathPrefix("/avatar").Handler(avaRoutes) // add avatar handlers

	authRouter.HandleFunc("/menus", server.newMenu).Methods("POST")
	authRouter.HandleFunc("/menus", server.allMenus).Methods("GET")
	authRouter.HandleFunc("/menus/{uuid}", server.getMenu).Methods("GET")
	authRouter.HandleFunc("/menus/{uuid}", server.updateMenu).Methods("PUT")
	authRouter.HandleFunc("/menus/{uuid}", server.deleteMenu).Methods("DELETE")

	authRouter.HandleFunc("/orders", server.newOrder).Methods("POST")
	authRouter.HandleFunc("/orders", server.allOrders).Methods("GET")
	authRouter.HandleFunc("/orders/{uuid}", server.getOrder).Methods("GET")
	authRouter.HandleFunc("/orders/{uuid}", server.updateOrder).Methods("PUT")
	authRouter.HandleFunc("/orders/{uuid}", server.deleteOrder).Methods("DELETE")

	router.HandleFunc("/users", server.registerUser).Methods("POST")
	authRouter.HandleFunc("/users", server.allUsers).Methods("GET")
	authRouter.HandleFunc("/users/{uuid}", server.getUser).Methods("GET")
	authRouter.HandleFunc("/users/{uuid}", server.updateUser).Methods("PUT")
	authRouter.HandleFunc("/users/{uuid}", server.deleteUser).Methods("DELETE")

	srv := &http.Server{
		Handler: router,
		Addr:    server.ctx.Value(AddressKey).(string),
	}
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Ctx(server.ctx).Fatal().Err(err).Msg("listen and serve crashed")
		}
	}()
	for {
		select {
		case <-server.ctx.Done():
			srv.Shutdown(server.ctx)
			log.Ctx(server.ctx).Debug().Msg("HTTP server shutdown complete")
			return
		}
	}
}
