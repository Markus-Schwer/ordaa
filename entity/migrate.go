package entity

import (
	"context"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/rs/zerolog/log"
)

func Migrate(ctx context.Context, databaseUrl string) error {
	log.Ctx(ctx).Info().Msg("running database migrations")

	m, err := migrate.New("file://db/migrations", databaseUrl)
	if err != nil {
		return fmt.Errorf("database migration preparation failed: %w", err)
	}
	if err := m.Up(); err != nil {
		if err != migrate.ErrNoChange {
			return fmt.Errorf("database migration execution failed: %w", err)
		} else {
			log.Ctx(ctx).Info().Msg("executed database migrations: no change")
		}
	}

	log.Ctx(ctx).Info().Msg("successfully ran database migrations")

	return nil
}
