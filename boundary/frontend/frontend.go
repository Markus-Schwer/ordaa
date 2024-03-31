package frontend

import (
	"context"

	"github.com/a-h/templ"
	"github.com/gorilla/mux"
	"gitlab.com/sfz.aalen/hackwerk/dotinder/entity"
)

const AddressKey = "ADDRESS"

type FrontendBoundary struct {
	ctx  context.Context
	repo *entity.Repository
}

func NewFrontendBoundary(ctx context.Context, repo *entity.Repository) *FrontendBoundary {
	return &FrontendBoundary{ctx: ctx, repo: repo}
}

func (server *FrontendBoundary) Start(router *mux.Router, authRouter *mux.Router) {
	router.Handle("/", templ.Handler(index()))
	router.HandleFunc("/orders", server.allOrders)
	router.HandleFunc("/orders/{uuid}", server.getOrder)
	router.HandleFunc("/menus", server.allMenus)
	router.HandleFunc("/menus/{uuid}", server.getMenu)
	router.Handle("/admin", templ.Handler(admin()))
}


