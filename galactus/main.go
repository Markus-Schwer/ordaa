package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"syscall"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gitlab.com/sfz.aalen/hackwerk/dotinder/galactus/orders"
)

const (
	AddressKey      = "Address"
	OmegaStarUrlKey = "OmegaStarURL"
)

func main() {
	var address, omegaStarUrl string
	var verbose bool
	flag.BoolVar(&verbose, "v", false, "verbose output: sets the log level to debug")
	flag.StringVar(&address, "address", "0.0.0.0:80", "the address including port of the service")
    flag.StringVar(&omegaStarUrl, "omega-star", "", "the URL of the omega star service")
	flag.Parse()

	ctx, cancel := context.WithCancel(context.Background())
	ctx = context.WithValue(ctx, AddressKey, address)
	ctx = context.WithValue(ctx, OmegaStarUrlKey, omegaStarUrl)
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	ctx = log.With().Str("service", "galactus").Logger().WithContext(ctx)
	if verbose {
		log.Ctx(ctx).Level(zerolog.DebugLevel)
	} else {
		log.Ctx(ctx).Level(zerolog.InfoLevel)
	}

	actionChan := make(chan orders.OrderAction, 10)
	mo, responseChan := orders.NewMultiOrders(ctx)
	restResponseChan := make(chan orders.OrderActionResponse, 10)
	// queueResponseChan := make(chan orders.OrderActionResponse, 10)

	go fanToAll(ctx, responseChan, restResponseChan)

	go mo.Start(actionChan)
	rest := NewRestInterface(ctx, actionChan, restResponseChan)
	go rest.start()
	// queue := NewQueueClient(ctx, actionChan, queueResponseChan)
	// go queue.start()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c
	log.Ctx(ctx).Info().Msg("received shutdown signal")
	cancel()
	log.Ctx(ctx).Info().Msg("finished graceful shutdown")
	os.Exit(0)
}

func fanToAll(ctx context.Context, src chan orders.OrderActionResponse, out ...chan orders.OrderActionResponse) {
	for {
		select {
		case <-ctx.Done():
			return
		case res := <-src:
			for _, outChan := range out {
				outChan <- res
			}
		}
	}
}
