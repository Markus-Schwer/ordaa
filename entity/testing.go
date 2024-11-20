package entity

import (
	"context"

	"github.com/glebarez/sqlite"
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

