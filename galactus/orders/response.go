package orders

import (
	"fmt"

	"github.com/google/uuid"
)

type OrderActionResponse interface {
	Uuid() uuid.UUID
	Error() error
}

type NoActiveOrder struct {
	uuid    uuid.UUID
	orderNo int
}

func (response *NoActiveOrder) Uuid() uuid.UUID {
	return response.uuid
}

func (response *NoActiveOrder) Error() error {
	return fmt.Errorf("there is no active order with id '%d'", response.orderNo)
}

type Ok struct {
	uuid    uuid.UUID
	OrderNo int
}

func (response *Ok) Uuid() uuid.UUID {
	return response.uuid
}

func (response *Ok) Error() error {
	return nil
}

type OkWithOrderNo struct {
	uuid    uuid.UUID
	OrderNo int
}

func (response *OkWithOrderNo) Uuid() uuid.UUID {
	return response.uuid
}

func (response *OkWithOrderNo) Error() error {
	return nil
}

type OkWithOrder struct {
	uuid  uuid.UUID
	Order *Order
}

func (response *OkWithOrder) Uuid() uuid.UUID {
	return response.uuid
}

func (response *OkWithOrder) Error() error {
	return nil
}

type GenericError struct {
	uuid uuid.UUID
	err  error
}

func (response *GenericError) Uuid() uuid.UUID {
	return response.uuid
}

func (response *GenericError) Error() error {
	return response.err
}

type OkWithActiveOrders struct {
	uuid         uuid.UUID
	ActiveOrders []ActiveOrder
}

func (response *OkWithActiveOrders) Uuid() uuid.UUID {
	return response.uuid
}

func (response *OkWithActiveOrders) Error() error {
	return nil
}
