package entity

import (
	"fmt"
	"time"

	"github.com/gofrs/uuid"
	"gorm.io/gorm"
)

type OrderState = string

const (
	Open      = OrderState("open")
	Finalized = OrderState("finalized")
	Ordered   = OrderState("ordered")
	Delivered = OrderState("delivered")
)

type Order struct {
	UUID          *uuid.UUID `gorm:"column:uuid;primaryKey" json:"uuid"`
	Initiator     *uuid.UUID `gorm:"column:initiator" json:"initiator"`
	SugarPerson   *uuid.UUID `gorm:"column:sugar_person" json:"sugar_person"`
	State         OrderState `gorm:"column:state" json:"state" validate:"omitempty,oneof=open finalized ordered delivered"`
	OrderDeadline *time.Time `gorm:"column:order_deadline" json:"order_deadline"`
	Eta           *time.Time `gorm:"column:eta" json:"eta"`
	MenuUUID      *uuid.UUID `gorm:"column:menu_uuid" json:"menu_uuid"`
}

type OrderItem struct {
	UUID         *uuid.UUID `gorm:"column:uuid;primaryKey" json:"uuid"`
	Price        int        `gorm:"column:price" json:"price"`
	Paid         bool       `gorm:"column:paid" json:"paid"`
	User         *uuid.UUID `gorm:"column:order_user" json:"order_user"`
	OrderUUID    *uuid.UUID `gorm:"column:order_uuid" json:"order_uuid"`
	MenuItemUUID *uuid.UUID `gorm:"column:menu_item_uuid" json:"menu_item_uuid" validate:"required"`
}

func (order *Order) BeforeCreate(tx *gorm.DB) (err error) {
	newUUID, err := uuid.NewV4()
	if err != nil {
		return fmt.Errorf("%w: %w", ErrCannotCreatUUID, err)
	}

	order.UUID = &newUUID

	return nil
}

func (orderItem *OrderItem) BeforeCreate(tx *gorm.DB) (err error) {
	newUUID, err := uuid.NewV4()
	if err != nil {
		return fmt.Errorf("%w: %w", ErrCannotCreatUUID, err)
	}

	orderItem.UUID = &newUUID

	return nil
}
