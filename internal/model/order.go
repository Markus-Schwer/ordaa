package model

import (
	"time"

	"github.com/gofrs/uuid"
)

type OrderState = string

const (
	Open      = OrderState("open")
	Finalized = OrderState("finalized")
	Ordered   = OrderState("ordered")
	Delivered = OrderState("delivered")
)

type Order struct {
	UUID          *uuid.UUID
	Initiator     *uuid.UUID
	SugarPerson   *uuid.UUID
	State         OrderState
	OrderDeadline *time.Time
	Eta           *time.Time
	MenuUUID      *uuid.UUID
}

type OrderItem struct {
	UUID         *uuid.UUID
	Price        int
	Paid         bool
	User         *uuid.UUID
	OrderUUID    *uuid.UUID
	MenuItemUUID *uuid.UUID
}
