package main

import (
	"context"
	"encoding/json"
	"flag"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gitlab.com/sfz.aalen/hackwerk/dotinder/boundary"
	"gitlab.com/sfz.aalen/hackwerk/dotinder/boundary/auth"
	"gitlab.com/sfz.aalen/hackwerk/dotinder/boundary/matrix"
	"gitlab.com/sfz.aalen/hackwerk/dotinder/boundary/rest"
	"gitlab.com/sfz.aalen/hackwerk/dotinder/boundary/tui"
	"gitlab.com/sfz.aalen/hackwerk/dotinder/entity"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	var verbose, jsonFormat bool
	flag.BoolVar(&verbose, "v", false, "verbose output: sets the log level to debug")
	flag.BoolVar(&jsonFormat, "j", false, "logging in json format")
	flag.Parse()

	ctx, cancel := context.WithCancel(context.Background())
	if !jsonFormat {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}

	ctx = log.With().Str("service", "dotinder").Logger().WithContext(ctx)
	if verbose {
		log.Ctx(ctx).Level(zerolog.DebugLevel)
	} else {
		log.Ctx(ctx).Level(zerolog.InfoLevel)
	}

	err := godotenv.Load()
	if err != nil {
		log.Ctx(ctx).Fatal().Err(err).Msg("Error loading .env file")
	}

	databaseUrl := os.Getenv(entity.DatabaseUrlKey)
	address := os.Getenv(boundary.AddressKey)
	homeserver := os.Getenv(matrix.HomeserverUrlKey)
	matrixUsername := os.Getenv(matrix.MatrixUsernameKey)
	matrixPassword := os.Getenv(matrix.MatrixPasswordKey)
	matrixRooms := strings.Split(os.Getenv(matrix.MatrixRoomsKey), ",")

	ctx = context.WithValue(ctx, entity.DatabaseUrlKey, databaseUrl)
	ctx = context.WithValue(ctx, boundary.AddressKey, address)
	ctx = context.WithValue(ctx, matrix.HomeserverUrlKey, homeserver)
	ctx = context.WithValue(ctx, matrix.MatrixUsernameKey, matrixUsername)
	ctx = context.WithValue(ctx, matrix.MatrixPasswordKey, matrixPassword)
	ctx = context.WithValue(ctx, matrix.MatrixRoomsKey, matrixRooms)

	if err := entity.Migrate(ctx, databaseUrl); err != nil {
		log.Ctx(ctx).Fatal().Err(err).Msg("database migration failed")
	}

	repo, err := entity.NewRepository(ctx, databaseUrl)
	if err != nil {
		log.Ctx(ctx).Fatal().Err(err).Msg(err.Error())
	}

	router := echo.New()
	authService := auth.NewAuthService(ctx, repo)

	rest.NewRestBoundary(ctx, repo, authService).Start(router)
	//frontend.NewFrontendBoundary(ctx, repo, authService).Start(router)
	go matrix.NewMatrixBoundary(ctx, repo).Start()

	go boundary.StartHttpServer(ctx, router)

	tuiServer := tui.NewSshTuiServer(ctx, repo)
	go tuiServer.Start()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c
	log.Ctx(ctx).Info().Msg("received shutdown signal")
	cancel()
	log.Ctx(ctx).Info().Msg("finished graceful shutdown")
	os.Exit(0)
}

func importMenu(filename string) (*entity.Menu, error) {
	bytes, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var menu entity.Menu
	if err := json.Unmarshal(bytes, &menu); err != nil {
		return nil, err
	}

	return &menu, nil
}
