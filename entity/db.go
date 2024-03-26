package entity

import (
	"context"

	"github.com/jmoiron/sqlx"
)

const DatabaseUrlKey = "DATABASE_URL"

type Repository struct {
	Pool sqlx.DB
	ctx  context.Context
}

func NewRepository(ctx context.Context, databaseUrl string) (*Repository, error) {
	db, err := sqlx.Connect("postgres", databaseUrl)
	if err != nil {
		return nil, err
	}

	return &Repository{
		Pool: *db,
		ctx:  ctx,
	}, nil
}
