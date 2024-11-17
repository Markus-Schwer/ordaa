package entity

import (
	"fmt"
	"github.com/gofrs/uuid"
	"gorm.io/gorm"
)

type Menu struct {
	Uuid  *uuid.UUID `gorm:"column:uuid;primaryKey" json:"uuid"`
	Name  string     `gorm:"column:name" json:"name" validate:"required"`
	Url   string     `gorm:"column:url" json:"url"`
	Items []MenuItem `gorm:"foreignKey:menu_uuid" json:"items"`
}

type MenuItem struct {
	Uuid      *uuid.UUID `gorm:"column:uuid;primaryKey" json:"uuid"`
	ShortName string     `gorm:"column:short_name" json:"short_name" validate:"required"`
	Name      string     `gorm:"column:name" json:"name" validate:"required"`
	Price     int        `gorm:"column:price" json:"price" validate:"required"`
	MenuUuid  *uuid.UUID  `gorm:"column:menu_uuid" json:"menu_uuid" validate:"required"`
}

func (menu *Menu) BeforeCreate(tx *gorm.DB) (err error) {
	newUuid, err := uuid.NewV4()
	if err != nil {
		return fmt.Errorf("could not create uuid: %w", err)
	}

	menu.Uuid = &newUuid
	return nil
}

func (menuItem *MenuItem) BeforeCreate(tx *gorm.DB) (err error) {
	newUuid, err := uuid.NewV4()
	if err != nil {
		return fmt.Errorf("could not create uuid: %w", err)
	}

	menuItem.Uuid = &newUuid
	return nil
}

func (*RepositoryImpl) GetAllMenus(tx *gorm.DB) ([]Menu, error) {
	menus := []Menu{}
	err := tx.Model(&Menu{}).Preload("Items").Find(&menus).Error
	if err != nil {
		return nil, fmt.Errorf("could not get all menus from db: %w", err)
	}

	return menus, nil
}

func (repo *RepositoryImpl) GetMenu(tx *gorm.DB, menuUuid *uuid.UUID) (*Menu, error) {
	var menu Menu
	if err := tx.Model(&Menu{}).Preload("Items").First(&menu, menuUuid).Error; err != nil {
		return nil, fmt.Errorf("failed to get menu %s: %w", menuUuid, err)
	}

	return &menu, nil
}

func (repo *RepositoryImpl) GetMenuByName(tx *gorm.DB, name string) (*Menu, error) {
	var menu Menu
	if err := tx.Model(&Menu{}).Preload("Items").Where(&Menu{Name: name}).First(&menu).Error; err != nil {
		return nil, fmt.Errorf("failed to get menu %s: %w", name, err)
	}

	return &menu, nil
}

func (repo *RepositoryImpl) GetMenuItem(tx *gorm.DB, menuItemUuid *uuid.UUID) (*MenuItem, error) {
	var menuItem MenuItem
	if err := tx.First(&MenuItem{}, menuItemUuid).Error; err != nil {
		return nil, fmt.Errorf("failed to get menu item %s: %w", menuItemUuid, err)
	}

	return &menuItem, nil
}

func (repo *RepositoryImpl) CreateMenu(tx *gorm.DB, menu *Menu) (*Menu, error) {
	err := tx.Create(&menu).Error
	if err != nil {
		return nil, fmt.Errorf("could not create menu %s: %w", menu.Name, err)
	}

	return menu, nil
}

func (repo *RepositoryImpl) UpdateMenu(tx *gorm.DB, menuUuid *uuid.UUID, menu *Menu) (*Menu, error) {
	existingMenu, err := repo.GetMenu(tx, menuUuid)
	if err != nil {
		return nil, err
	}

	existingMenu.Name = menu.Name
	existingMenu.Url = menu.Url
	existingMenu.Items = menu.Items
	err = tx.Save(existingMenu).Error
	if err != nil {
		return nil, fmt.Errorf("could not update menu %s: %w", menuUuid, err)
	}

	return existingMenu, nil
}

func (repo *RepositoryImpl) CreateMenuItem(tx *gorm.DB, menuItem *MenuItem, menuUuid *uuid.UUID) (*MenuItem, error) {
	err := tx.Create(&menuItem).Error
	if err != nil {
		return nil, fmt.Errorf("could not create menu item %s: %w", menuItem.ShortName, err)
	}

	return menuItem, nil
}

func (repo *RepositoryImpl) DeleteMenuItem(tx *gorm.DB, menuItemUuid *uuid.UUID) error {
	err := tx.Delete(&MenuItem{}, menuItemUuid).Error
	if err != nil {
		return fmt.Errorf("could not delete menu item %s: %w", menuItemUuid, err)
	}

	return nil
}

func (repo *RepositoryImpl) DeleteMenu(tx *gorm.DB, menuUuid *uuid.UUID) error {
	err := tx.Delete(&Menu{}, menuUuid).Error
	if err != nil {
		return fmt.Errorf("could not delete menu %s: %w", menuUuid, err)
	}

	return nil
}
