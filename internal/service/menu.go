package service

import (
	"context"

	"github.com/gofrs/uuid"

	"github.com/Markus-Schwer/ordaa/internal/entity"
)

type MenuRepository interface {
	GetAllMenus(ctx context.Context) ([]entity.Menu, error)
	GetMenu(ctx context.Context, menuUUID *uuid.UUID) (*entity.Menu, error)
	GetMenuByName(ctx context.Context, name string) (*entity.Menu, error)
	GetMenuItem(ctx context.Context, menuItemUUID *uuid.UUID) (*entity.MenuItem, error)
	GetMenuItemByShortName(ctx context.Context, menuUUID *uuid.UUID, shortName string) (*entity.MenuItem, error)
	CreateMenu(ctx context.Context, menu *entity.Menu) (*entity.Menu, error)
	UpdateMenu(ctx context.Context, menuUUID *uuid.UUID, menu *entity.Menu) (*entity.Menu, error)
	CreateMenuItem(ctx context.Context, menuItem *entity.MenuItem) (*entity.MenuItem, error)
	DeleteMenuItem(ctx context.Context, menuItemUUID *uuid.UUID) error
	DeleteMenu(ctx context.Context, menuUUID *uuid.UUID) error
}

type MenuService struct {
	MenuRepository MenuRepository
}

func (s *MenuService) GetAllMenus(ctx context.Context) ([]entity.Menu, error) {
	return s.MenuRepository.GetAllMenus(ctx)
}

func (s *MenuService) GetMenu(ctx context.Context, uuid *uuid.UUID) (*entity.Menu, error) {
	return s.MenuRepository.GetMenu(ctx, uuid)
}

func (s *MenuService) CreateMenu(ctx context.Context, user *entity.Menu) (*entity.Menu, error) {
	return s.MenuRepository.CreateMenu(ctx, user)
}

func (s *MenuService) UpdateMenu(ctx context.Context, uuid *uuid.UUID, user *entity.Menu) (*entity.Menu, error) {
	return s.MenuRepository.UpdateMenu(ctx, uuid, user)
}

func (s *MenuService) DeleteMenu(ctx context.Context, uuid *uuid.UUID) error {
	return s.MenuRepository.DeleteMenu(ctx, uuid)
}

func (s *MenuService) ImportMenu(ctx context.Context, menu *entity.Menu) error {
	if _, err := s.MenuRepository.CreateMenu(ctx, menu); err != nil {
		return err
	}

	return nil
}
