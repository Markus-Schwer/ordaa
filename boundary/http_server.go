package boundary

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

const AddressKey = "ADDRESS"

func StartHttpServer(ctx context.Context, router *echo.Echo) {
	srv := &http.Server{
		Handler: router,
		Addr:    ctx.Value(AddressKey).(string),
	}
	go func() {
		log.Ctx(ctx).Info().Msg("starting http server")
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
