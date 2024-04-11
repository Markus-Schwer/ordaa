package frontend

import (
	"context"

	"github.com/gorilla/mux"
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

func (server *FrontendBoundary) Start(router *mux.Router, authRouter *mux.Router) {
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
