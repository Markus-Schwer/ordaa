package entity

import (
	"context"
	"errors"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres" // add postgres support for migrations
	_ "github.com/golang-migrate/migrate/v4/source/file"       // add files as a source for migrations
	"github.com/rs/zerolog/log"
)

func Migrate(ctx context.Context, databaseURL string) error {
	log.Ctx(ctx).Info().Msg("running database migrations")

	m, err := migrate.New("file://db/migrations", databaseURL)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrMigrationPreparationFailed, err)
	}

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("%w: %w", ErrMigrationFailed, err)
	}

	log.Ctx(ctx).Info().Msg("executed database migrations: no change")
	log.Ctx(ctx).Info().Msg("successfully ran database migrations")

	return nil
}
