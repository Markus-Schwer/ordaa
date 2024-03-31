package rest

import (
	"context"
	"github.com/gorilla/mux"
	"gitlab.com/sfz.aalen/hackwerk/dotinder/entity"
)

type RestBoundary struct {
	ctx  context.Context
	repo *entity.Repository
}

func NewRestBoundary(ctx context.Context, repo *entity.Repository) *RestBoundary {
	return &RestBoundary{ctx: ctx, repo: repo}
}

func (server *RestBoundary) Start(router *mux.Router, authRouter *mux.Router) {
	authRouter.HandleFunc("/api/menus", server.newMenu).Methods("POST")
	authRouter.HandleFunc("/api/menus", server.allMenus).Methods("GET")
	authRouter.HandleFunc("/api/menus/{uuid}", server.getMenu).Methods("GET")
	authRouter.HandleFunc("/api/menus/{uuid}", server.updateMenu).Methods("PUT")
	authRouter.HandleFunc("/api/menus/{uuid}", server.deleteMenu).Methods("DELETE")

	authRouter.HandleFunc("/api/orders", server.newOrder).Methods("POST")
	authRouter.HandleFunc("/api/orders", server.allOrders).Methods("GET")
	authRouter.HandleFunc("/api/orders/{uuid}", server.getOrder).Methods("GET")
	authRouter.HandleFunc("/api/orders/{uuid}", server.updateOrder).Methods("PUT")
	authRouter.HandleFunc("/api/orders/{uuid}", server.deleteOrder).Methods("DELETE")

	router.HandleFunc("/api/users", server.registerUser).Methods("POST")
	authRouter.HandleFunc("/api/users", server.allUsers).Methods("GET")
	authRouter.HandleFunc("/api/users/{uuid}", server.getUser).Methods("GET")
	authRouter.HandleFunc("/api/users/{uuid}", server.updateUser).Methods("PUT")
	authRouter.HandleFunc("/api/users/{uuid}", server.deleteUser).Methods("DELETE")
}
