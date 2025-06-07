package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
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
	"github.com/Markus-Schwer/ordaa/internal/config"
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

	logConfig, err := config.LoadLogConfig()
	if err != nil {
		return err
	}

	dbConfig, err := config.LoadDatabaseConfig()
	if err != nil {
		return err
	}

	matrixConfig, err := config.LoadMatrixConfig()
	if err != nil {
		return err
	}

	if !logConfig.JSON {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}

	ctx = log.With().Str("service", "ordaa").Logger().WithContext(ctx)

	level, err := zerolog.ParseLevel(logConfig.Level)
	if err != nil {
		return err
	}

	zerolog.SetGlobalLevel(level)

	log.Ctx(ctx).Info().Msg("starting ordaa")

	if err := entity.Migrate(ctx, dbConfig.URL); err != nil {
		return fmt.Errorf("database migration failed: %w", err)
	}

	db, err := gorm.Open(postgres.Open(dbConfig.URL), &gorm.Config{
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

	matrixBoundary, err := matrix.NewMatrixBoundary(ctx, matrixConfig, userService, orderService)
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
