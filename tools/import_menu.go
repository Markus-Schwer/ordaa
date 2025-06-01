package main

import (
	"context"
	"encoding/json"
	"flag"
	"os"

	_ "github.com/joho/godotenv/autoload"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/Markus-Schwer/ordaa/internal/entity"
	"github.com/Markus-Schwer/ordaa/internal/repository"
	"github.com/Markus-Schwer/ordaa/internal/service"
)

func main() {
	ctx := context.Background()

	var file string

	flag.StringVar(&file, "f", "sangam.json", "file to import")
	flag.Parse()

	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	zerolog.SetGlobalLevel(zerolog.InfoLevel)

	databaseURL := os.Getenv("DATABASE_URL")

	if err := entity.Migrate(ctx, databaseURL); err != nil {
		log.Fatal().Err(err).Msg("database migration failed")
	}

	db, err := gorm.Open(postgres.Open(databaseURL), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		log.Fatal().Err(err).Msg("connecting to database")
	}

	menuRepository := &repository.MenuRepository{DB: db}
	menuService := &service.MenuService{MenuRepository: menuRepository}

	menu, err := importMenu(file)
	if err != nil {
		log.Fatal().Err(err).Msg("reading menu json")
	}

	if err = menuService.ImportMenu(ctx, menu); err != nil {
		log.Fatal().Err(err).Msg("importing menu")
	}

	log.Info().Msg("successfully imported menu")
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
