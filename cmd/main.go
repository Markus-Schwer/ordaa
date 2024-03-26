package main

import (
	"context"
	"encoding/json"
	"flag"
	"io/ioutil"
	"os"
	"os/signal"
	"syscall"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gitlab.com/sfz.aalen/hackwerk/dotinder/boundary/rest"
	"gitlab.com/sfz.aalen/hackwerk/dotinder/entity"

	_ "github.com/lib/pq"
)

func main() {
	var verbose, jsonFormat bool
	flag.BoolVar(&verbose, "v", false, "verbose output: sets the log level to debug")
	flag.BoolVar(&jsonFormat, "j", false, "logging in json format")
	flag.Parse()
	databaseUrl := os.Getenv(entity.DatabaseUrlKey)
	address := os.Getenv(rest.AddressKey)

	ctx, cancel := context.WithCancel(context.Background())
	ctx = context.WithValue(ctx, entity.DatabaseUrlKey, databaseUrl)
	ctx = context.WithValue(ctx, rest.AddressKey, address)
	if !jsonFormat {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}
	ctx = log.With().Str("service", "dotinder").Logger().WithContext(ctx)
	if verbose {
		log.Ctx(ctx).Level(zerolog.DebugLevel)
	} else {
		log.Ctx(ctx).Level(zerolog.InfoLevel)
	}

	if err := entity.Migrate(ctx, databaseUrl); err != nil {
		log.Ctx(ctx).Fatal().Err(err).Msg("database migration failed")
	}

	repo, err := entity.NewRepository(ctx, databaseUrl)
	if err != nil {
		log.Ctx(ctx).Fatal().Err(err).Msg(err.Error())
	}

	server := rest.NewRestBoundary(ctx, repo)
	go server.Start()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c
	log.Ctx(ctx).Info().Msg("received shutdown signal")
	cancel()
	log.Ctx(ctx).Info().Msg("finished graceful shutdown")
	os.Exit(0)
}

func importMenu(filename string) (*entity.NewMenu, error) {
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var menu entity.NewMenu
	if err := json.Unmarshal(bytes, &menu); err != nil {
		return nil, err
	}

	return &menu, nil
}
