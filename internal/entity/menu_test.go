package entity

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

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
		menuItems := []MenuItem{}
		menu := &Menu{Name: "Test Menu", Url: "https://sangam-aalen.de", Items: menuItems}
		createdMenu, err := repo.CreateMenu(tx, menu)
		if err != nil {
			t.Fatal(err)
		}

		foundMenu, err := repo.GetMenu(tx, createdMenu.Uuid)
		if err != nil {
			return nil
		}

		assert.Equal(t, Menu{Uuid: foundMenu.Uuid, Name: menu.Name, Url: menu.Url, Items: menuItems}, *foundMenu)
		return nil
	}))
}
