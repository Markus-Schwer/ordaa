package entity

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
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

func TestGetAllMenus(t *testing.T) {
	ctx := context.Background()
	repo, err := NewTestRepository(ctx)
	if err != nil {
		t.Fatal(err)
	}

	assert.NoError(t, repo.Transaction(func(tx *gorm.DB) error {
		menus, err := repo.GetAllMenus(tx)
		if err != nil {
			return nil
		}

		assert.Equal(t, []Menu{}, menus)
		return nil
	}))
}

func TestGetMenu(t *testing.T) {
	ctx := context.Background()
	repo, err := NewTestRepository(ctx)
	if err != nil {
		t.Fatal(err)
	}

	assert.NoError(t, repo.Transaction(func(tx *gorm.DB) error {
		menus, err := repo.GetMenu(tx, )
		if err != nil {
			return nil
		}

		assert.Equal(t, []Menu{}, menus)
		return nil
	}))
}
