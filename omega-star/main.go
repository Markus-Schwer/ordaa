package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"syscall"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

const (
	AddressKey = "Address"
)

func main() {
	var verbose, jsonFormat bool
	flag.BoolVar(&verbose, "v", false, "verbose output: sets the log level to debug")
	flag.BoolVar(&jsonFormat, "j", false, "logging in json format")
	address := os.Getenv("ADDRESS")
	flag.Parse()

	ctx, cancel := context.WithCancel(context.Background())
	ctx = context.WithValue(ctx, AddressKey, address)
	if !jsonFormat {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}
	ctx = log.With().Str("service", "omega-star").Logger().WithContext(ctx)
	if verbose {
		log.Ctx(ctx).Level(zerolog.DebugLevel)
	} else {
		log.Ctx(ctx).Level(zerolog.InfoLevel)
	}

	server := NewMenuServer(ctx)
	log.Ctx(ctx).Info().Msgf("starting omega star at '%s'", address)
	go server.start()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c
	log.Ctx(ctx).Info().Msg("received shutdown signal")
	cancel()
	log.Ctx(ctx).Info().Msg("finished graceful shutdown")
	os.Exit(0)
}
