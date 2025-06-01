package entity

import (
	"fmt"

	"github.com/gofrs/uuid"
	"gorm.io/gorm"
)

type Menu struct {
	UUID  *uuid.UUID `gorm:"column:uuid;primaryKey" json:"uuid"`
	Name  string     `gorm:"column:name" json:"name" validate:"required"`
	URL   string     `gorm:"column:url" json:"url"`
	Items []MenuItem `gorm:"foreignKey:menu_uuid" json:"items"`
}

type MenuItem struct {
	UUID      *uuid.UUID `gorm:"column:uuid;primaryKey" json:"uuid"`
	ShortName string     `gorm:"column:short_name" json:"short_name" validate:"required"`
	Name      string     `gorm:"column:name" json:"name" validate:"required"`
	Price     int        `gorm:"column:price" json:"price" validate:"required"`
	MenuUUID  *uuid.UUID `gorm:"column:menu_uuid" json:"menu_uuid" validate:"required"`
}

func (menu *Menu) BeforeCreate(tx *gorm.DB) (err error) {
	newUUID, err := uuid.NewV4()
	if err != nil {
		return fmt.Errorf("%w: %w", ErrCannotCreatUUID, err)
	}

	menu.UUID = &newUUID

	return nil
}

func (menuItem *MenuItem) BeforeCreate(tx *gorm.DB) (err error) {
	newUUID, err := uuid.NewV4()
	if err != nil {
		return fmt.Errorf("%w: %w", ErrCannotCreatUUID, err)
	}

	menuItem.UUID = &newUUID

	return nil
}
