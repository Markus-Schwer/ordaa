package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/gofrs/uuid"
	"gorm.io/gorm"

	"github.com/Markus-Schwer/ordaa/internal/entity"
)

var (
	ErrCannotGetAllMenus = errors.New("could not get all menus from db")
	ErrGettingMenu       = errors.New("failed to get menu")
	ErrMenuNotFound      = errors.New("menu not found")
	ErrGettingMenuItem   = errors.New("failed to get menu item")
	ErrMenuItemNotFound  = errors.New("menu item not found")
	ErrCreatingMenu      = errors.New("could not create menu")
	ErrUpdatingMenu      = errors.New("could not update menu")
	ErrCreatingMenuItem  = errors.New("could not create menu item")
	ErrDeletingMenuItem  = errors.New("could not delete menu item")
	ErrDeletingMenu      = errors.New("could not delete menu")
)

type MenuRepository struct {
	DB *gorm.DB
}

func (r *MenuRepository) GetAllMenus(ctx context.Context) ([]entity.Menu, error) {
	menus := []entity.Menu{}

	err := r.DB.Model(&entity.Menu{}).Preload("Items").Find(&menus).Error
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrCannotGetAllMenus, err)
	}

	return menus, nil
}

func (r *MenuRepository) GetMenu(ctx context.Context, menuUUID *uuid.UUID) (*entity.Menu, error) {
	var menu entity.Menu

	err := r.DB.Model(&entity.Menu{}).Preload("Items").First(&menu, menuUUID).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("%w: %w", ErrMenuItemNotFound, err)
	} else if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrGettingMenu, err)
	}

	return &menu, nil
}

func (r *MenuRepository) GetMenuByName(ctx context.Context, name string) (*entity.Menu, error) {
	var menu entity.Menu

	err := r.DB.Model(&entity.Menu{}).Preload("Items").Where(&entity.Menu{Name: name}).First(&menu).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("%w: %w", ErrMenuNotFound, err)
	} else if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrGettingMenu, err)
	}

	return &menu, nil
}

func (r *MenuRepository) GetMenuItem(ctx context.Context, menuItemUUID *uuid.UUID) (*entity.MenuItem, error) {
	var menuItem entity.MenuItem

	err := r.DB.First(&menuItem, menuItemUUID).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("%w: %w", ErrMenuItemNotFound, err)
	} else if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrGettingMenuItem, err)
	}

	return &menuItem, nil
}

func (r *MenuRepository) GetMenuItemByShortName(ctx context.Context, menuUUID *uuid.UUID, shortName string) (*entity.MenuItem, error) {
	var menuItem entity.MenuItem

	err := r.DB.First(&menuItem, entity.MenuItem{MenuUUID: menuUUID, ShortName: shortName}).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("%w: %w", ErrMenuItemNotFound, err)
	} else if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrGettingMenuItem, err)
	}

	return &menuItem, nil
}

func (r *MenuRepository) CreateMenu(ctx context.Context, menu *entity.Menu) (*entity.Menu, error) {
	tx := r.DB.Begin()

	err := tx.Create(&menu).Error
	if err != nil {
		_ = tx.Rollback()
		return nil, fmt.Errorf("%w: %w", ErrCreatingMenu, err)
	}

	_ = tx.Commit()

	return menu, nil
}

func (r *MenuRepository) UpdateMenu(ctx context.Context, menuUUID *uuid.UUID, menu *entity.Menu) (*entity.Menu, error) {
	tx := r.DB.Begin()

	existingMenu, err := r.GetMenu(ctx, menuUUID)
	if err != nil {
		_ = tx.Rollback()
		return nil, err
	}

	existingMenu.Name = menu.Name
	existingMenu.URL = menu.URL
	existingMenu.Items = menu.Items

	err = tx.Save(existingMenu).Error
	if err != nil {
		_ = tx.Rollback()
		return nil, fmt.Errorf("%w: %w", ErrUpdatingMenu, err)
	}

	_ = tx.Commit()

	return existingMenu, nil
}

func (r *MenuRepository) CreateMenuItem(ctx context.Context, menuItem *entity.MenuItem) (*entity.MenuItem, error) {
	tx := r.DB.Begin()

	err := tx.Create(&menuItem).Error
	if err != nil {
		_ = tx.Rollback()
		return nil, fmt.Errorf("%w: %w", ErrCreatingMenuItem, err)
	}

	_ = tx.Commit()

	return menuItem, nil
}

func (r *MenuRepository) DeleteMenuItem(ctx context.Context, menuItemUUID *uuid.UUID) error {
	tx := r.DB.Begin()

	err := tx.Delete(&entity.MenuItem{}, menuItemUUID).Error
	if err != nil {
		_ = tx.Rollback()
		return fmt.Errorf("%w: %w", ErrDeletingMenuItem, err)
	}

	_ = tx.Commit()

	return nil
}

func (r *MenuRepository) DeleteMenu(ctx context.Context, menuUUID *uuid.UUID) error {
	tx := r.DB.Begin()

	err := tx.Delete(&entity.Menu{}, menuUUID).Error
	if err != nil {
		_ = tx.Rollback()
		return fmt.Errorf("%w: %w", ErrDeletingMenu, err)
	}

	_ = tx.Commit()

	return nil
}
