package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	_ "github.com/joho/godotenv/autoload"
	_ "github.com/lib/pq"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"golang.org/x/sync/errgroup"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/Markus-Schwer/ordaa/internal/boundary/matrix"
	"github.com/Markus-Schwer/ordaa/internal/entity"
	"github.com/Markus-Schwer/ordaa/internal/repository"
	"github.com/Markus-Schwer/ordaa/internal/service"
)

func main() {
	if err := Run(); err != nil {
		log.Fatal().Err(err).Msg("running ordaa")
	}
}

func Run() error {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	var verbose, jsonFormat bool

	flag.BoolVar(&verbose, "v", false, "verbose output: sets the log level to debug")
	flag.BoolVar(&jsonFormat, "j", false, "logging in json format")
	flag.Parse()

	if !jsonFormat {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}

	ctx = log.With().Str("service", "ordaa").Logger().WithContext(ctx)

	if verbose {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	} else {
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}

	databaseURL := os.Getenv("DATABASE_URL")
	homeserver := os.Getenv("MATRIX_HOMESERVER")
	matrixUsername := os.Getenv("MATRIX_USERNAME")
	matrixPassword := os.Getenv("MATRIX_PASSWORD")
	matrixRooms := strings.Split(os.Getenv("MATRIX_ROOMS"), ",")

	log.Ctx(ctx).Info().Msg("starting ordaa")

	if err := entity.Migrate(ctx, databaseURL); err != nil {
		return fmt.Errorf("database migration failed: %w", err)
	}

	db, err := gorm.Open(postgres.Open(databaseURL), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		return err
	}

	userRepository := &repository.UserRepository{DB: db}
	menuRepository := &repository.MenuRepository{DB: db}
	orderRepository := &repository.OrderRepository{DB: db, MenuRepository: *menuRepository}

	userService := &service.UserService{UserRepository: userRepository}
	orderService := &service.OrderService{OrderRepository: orderRepository, MenuRepository: menuRepository}

	g, gCtx := errgroup.WithContext(ctx)

	matrixBoundary, err := matrix.NewMatrixBoundary(ctx, homeserver, matrixUsername, matrixPassword, matrixRooms, userService, orderService)
	if err != nil {
		return err
	}

	g.Go(func() error {
		return matrixBoundary.Start(ctx)
	})

	g.Go(func() error {
		<-gCtx.Done()
		return matrixBoundary.Stop()
	})

	if err := g.Wait(); err != nil {
		return fmt.Errorf("shutting down: %w", err)
	}

	return nil
}
