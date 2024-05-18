package entity

import (
	"fmt"
	"time"

	"github.com/gofrs/uuid"
	"gorm.io/gorm"
)

type Order struct {
	Uuid          *uuid.UUID  `gorm:"column:uuid;primaryKey" json:"uuid"`
	Initiator     uuid.UUID   `gorm:"column:initiator" json:"initiator"`
	SugarPerson   *uuid.UUID  `gorm:"column:sugar_person" json:"sugar_person"`
	State         string      `gorm:"column:state" json:"state"`
	OrderDeadline *time.Time  `gorm:"column:order_deadline" json:"order_deadline"`
	Eta           *time.Time  `gorm:"column:eta" json:"eta"`
	MenuUuid      uuid.UUID   `gorm:"column:menu_uuid" json:"menu_uuid"`
	Items         []OrderItem `gorm:"foreignKey:order_uuid" json:"items"`
}

type OrderItem struct {
	Uuid         *uuid.UUID `gorm:"column:uuid;primaryKey" json:"uuid"`
	Price        int        `gorm:"column:price" json:"price"`
	Paid         bool       `gorm:"column:paid" json:"paid"`
	User         uuid.UUID  `gorm:"column:order_user" json:"order_user"`
	OrderUuid    uuid.UUID  `gorm:"column:order_uuid" json:"order_uuid"`
	MenuItemUuid uuid.UUID  `gorm:"column:menu_item_uuid" json:"menu_item_uuid"`
}

func (*Repository) GetAllOrders(tx *gorm.DB) ([]Order, error) {
	orders := []Order{}
	err := tx.Model(&Order{}).Preload("Items").Find(&orders).Error
	if err != nil {
		return nil, fmt.Errorf("could not get all orders from db: %w", err)
	}

	return orders, nil
}

func (*Repository) GetOrder(tx *gorm.DB, uuid uuid.UUID) (*Order, error) {
	var order Order
	if err := tx.Model(&Order{}).Preload("Items").First(&order, uuid).Error; err != nil {
		return nil, fmt.Errorf("error getting order %s: %w", uuid, err)
	}

	return &order, nil
}

func (*Repository) GetAllOrderItems(tx *gorm.DB, orderUuid uuid.UUID) ([]OrderItem, error) {
	orderItems := []OrderItem{}
	err := tx.Find(&orderItems).Error
	if err != nil {
		return nil, fmt.Errorf("could not get all order items from db: %w", err)
	}

	return orderItems, nil
}

func (*Repository) GetOrderItem(tx *gorm.DB, uuid uuid.UUID) (*OrderItem, error) {
	orderItem := OrderItem{}
	if err := tx.First(&orderItem, uuid).Error; err != nil {
		return nil, fmt.Errorf("error getting order item %s: %w", uuid, err)
	}

	return &orderItem, nil
}

func (*Repository) CreateOrderItem(tx *gorm.DB, order_uuid uuid.UUID, orderItem *OrderItem) (*OrderItem, error) {
	err := tx.Create(&orderItem).Error
	if err != nil {
		return nil, fmt.Errorf("could not create order item: %w", err)
	}

	return orderItem, nil
}

func (repo *Repository) CreateOrder(tx *gorm.DB, order *Order) (*Order, error) {
	err := tx.Create(&order).Error
	if err != nil {
		return nil, fmt.Errorf("could not create order: %w", err)
	}

	return order, nil
}

func (repo *Repository) UpdateOrder(tx *gorm.DB, orderUuid uuid.UUID, order *Order) (*Order, error) {
	existingOrder, err := repo.GetOrder(tx, orderUuid)
	if err != nil {
		return nil, fmt.Errorf("could not update order %s: %w", orderUuid, err)
	}

	existingOrder.Initiator = order.Initiator
	existingOrder.SugarPerson = order.SugarPerson
	existingOrder.State = order.State
	existingOrder.OrderDeadline = order.OrderDeadline
	existingOrder.Eta = order.Eta

	err = tx.Save(existingOrder).Error
	if err != nil {
		return nil, fmt.Errorf("could not update order %s: %w", orderUuid, err)
	}

	return existingOrder, nil
}

func (repo *Repository) UpdateOrderItem(tx *gorm.DB, orderItemUuid uuid.UUID, orderItem *OrderItem) (*OrderItem, error) {
	existingOrderItem, err := repo.GetOrderItem(tx, orderItemUuid)
	if err != nil {
		return nil, fmt.Errorf("could not update order item %s: %w", orderItemUuid, err)
	}

	existingOrderItem.User = orderItem.User
	existingOrderItem.Price = orderItem.Price
	existingOrderItem.Paid = orderItem.Paid
	existingOrderItem.MenuItemUuid = orderItem.MenuItemUuid

	err = tx.Save(existingOrderItem).Error
	if err != nil {
		return nil, fmt.Errorf("could not update order item %s: %w", orderItemUuid, err)
	}

	return existingOrderItem, nil
}

func (repo *Repository) DeleteOrderItem(tx *gorm.DB, orderItemUuid uuid.UUID) error {
	err := tx.Delete(&OrderItem{}, orderItemUuid).Error
	if err != nil {
		return fmt.Errorf("could not delete order item %s: %w", orderItemUuid, err)
	}

	return nil
}

func (repo *Repository) DeleteOrder(tx *gorm.DB, orderUuid uuid.UUID) error {
	err := tx.Delete(&OrderItem{}, orderUuid).Error
	if err != nil {
		return fmt.Errorf("could not delete order %s: %w", orderUuid, err)
	}

	return nil
}
