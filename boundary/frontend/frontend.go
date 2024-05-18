package frontend

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
	"gitlab.com/sfz.aalen/hackwerk/dotinder/boundary/auth"
	"gitlab.com/sfz.aalen/hackwerk/dotinder/entity"
)

const AddressKey = "ADDRESS"

type FrontendBoundary struct {
	ctx         context.Context
	repo        *entity.Repository
	authService *auth.AuthService
}

func NewFrontendBoundary(ctx context.Context, repo *entity.Repository, authService *auth.AuthService) *FrontendBoundary {
	return &FrontendBoundary{ctx: ctx, repo: repo, authService: authService}
}

func (server *FrontendBoundary) Start(router *echo.Echo) {
	authRouter := router.Group("/", auth.AuthMiddleware(server.authService, func(c echo.Context, err error) error {
		server.authService.Logout(c)
		return echo.NewHTTPError(http.StatusUnauthorized, "Unauthorized")
	}))

	authRouter := router.NewRoute().Subrouter()
	authRouter.Use(auth.AuthMiddleware(server.authService, mux.MiddlewareFunc(func(handler http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			server.authService.Logout(w, r)
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			handler.ServeHTTP(w, r)
		})
	})))

	router.HandleFunc("/", server.index)
	authRouter.HandleFunc("/orders", server.allOrders)
	authRouter.HandleFunc("/orders/{uuid}", server.getOrder)
	authRouter.HandleFunc("/menus", server.allMenus)
	authRouter.HandleFunc("/menus/{uuid}", server.getMenu)
	authRouter.HandleFunc("/admin", server.admin)
	router.HandleFunc("/login", server.login)
	router.HandleFunc("/logout", server.authService.Logout)
	router.HandleFunc("/signup", server.signup)
}
