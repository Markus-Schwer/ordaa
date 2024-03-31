package boundary

import (
	"context"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
)

const AddressKey = "ADDRESS"

func StartHttpServer(ctx context.Context, router *mux.Router) {
	srv := &http.Server{
		Handler: router,
		Addr:    ctx.Value(AddressKey).(string),
	}
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Ctx(ctx).Fatal().Err(err).Msg("listen and serve crashed")
		}
	}()
	for {
		select {
		case <-ctx.Done():
			srv.Shutdown(ctx)
			log.Ctx(ctx).Debug().Msg("HTTP server shutdown complete")
			return
		}
	}
}
