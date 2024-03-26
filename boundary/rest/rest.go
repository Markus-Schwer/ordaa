package rest

import (
	"context"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
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
	router := mux.NewRouter()

	router.HandleFunc("/menus", server.newMenu).Methods("POST")
	router.HandleFunc("/menus", server.allMenus).Methods("GET")
	router.HandleFunc("/menus/{uuid}", server.getMenu).Methods("GET")
	router.HandleFunc("/menus/{uuid}", server.updateMenu).Methods("PUT")
	router.HandleFunc("/menus/{uuid}", server.deleteMenu).Methods("DELETE")

	//router.HandleFunc("/orders", server.newOrder).Methods("POST")
	router.HandleFunc("/orders", server.allOrders).Methods("GET")
	//router.HandleFunc("/orders/{uuid}", server.getOrder).Methods("GET")
	//router.HandleFunc("/orders/{uuid}", server.updateOrder).Methods("PUT")
	//router.HandleFunc("/orders/{uuid}", server.deleteOrder).Methods("DELETE")

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
