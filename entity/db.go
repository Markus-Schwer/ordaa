package entity

import (
	"context"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const DatabaseUrlKey = "DATABASE_URL"

type Repository struct {
	Db  *gorm.DB
	ctx context.Context
}

func NewRepository(ctx context.Context, databaseUrl string) (*Repository, error) {
	db, err := gorm.Open(postgres.Open(databaseUrl), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return &Repository{
		Db:  db,
		ctx: ctx,
	}, nil
}
