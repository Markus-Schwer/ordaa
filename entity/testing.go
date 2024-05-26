package entity

import (
	"context"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func NewTestRepository(ctx context.Context) (*RepositoryImpl, error) {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	err = db.AutoMigrate(&Menu{}, &MenuItem{}, &Order{}, &OrderItem{}, &User{}, &PasswordUser{}, &MatrixUser{})
	if err != nil {
		return nil, err
	}

	return &RepositoryImpl{
		ctx: ctx,
		Db:  db,
	}, nil
}

