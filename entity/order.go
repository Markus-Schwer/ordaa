package entity

import (
	"errors"
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
	Uuid          *uuid.UUID  `gorm:"column:uuid;primaryKey" json:"uuid"`
	Initiator     *uuid.UUID  `gorm:"column:initiator" json:"initiator"`
	SugarPerson   *uuid.UUID  `gorm:"column:sugar_person" json:"sugar_person"`
	State         OrderState  `gorm:"column:state" json:"state" validate:"omitempty,oneof=open finalized ordered delivered"`
	OrderDeadline *time.Time  `gorm:"column:order_deadline" json:"order_deadline"`
	Eta           *time.Time  `gorm:"column:eta" json:"eta"`
	MenuUuid      *uuid.UUID  `gorm:"column:menu_uuid" json:"menu_uuid"`
	Items         []OrderItem `gorm:"foreignKey:order_uuid" json:"items"`
}

type OrderItem struct {
	Uuid         *uuid.UUID `gorm:"column:uuid;primaryKey" json:"uuid"`
	Price        int        `gorm:"column:price" json:"price"`
	Paid         bool       `gorm:"column:paid" json:"paid"`
	User         *uuid.UUID `gorm:"column:order_user" json:"order_user"`
	OrderUuid    *uuid.UUID `gorm:"column:order_uuid" json:"order_uuid"`
	MenuItemUuid *uuid.UUID `gorm:"column:menu_item_uuid" json:"menu_item_uuid" validate:"required"`
}

func (order *Order) BeforeCreate(tx *gorm.DB) (err error) {
	newUuid, err := uuid.NewV4()
	if err != nil {
		return fmt.Errorf("%w: %w", ErrCannotCreatUuid, err)
	}

	order.Uuid = &newUuid
	return nil
}

func (orderItem *OrderItem) BeforeCreate(tx *gorm.DB) (err error) {
	newUuid, err := uuid.NewV4()
	if err != nil {
		return fmt.Errorf("%w: %w", ErrCannotCreatUuid, err)
	}

	orderItem.Uuid = &newUuid
	return nil
}

func (*RepositoryImpl) GetAllOrders(tx *gorm.DB) ([]Order, error) {
	orders := []Order{}
	err := tx.Model(&Order{}).Preload("Items").Find(&orders).Error
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrCannotGetAllOrder, err)
	}

	return orders, nil
}

func (*RepositoryImpl) GetOrder(tx *gorm.DB, uuid *uuid.UUID) (*Order, error) {
	var order Order
	err := tx.Model(&Order{}).Preload("Items").First(&order, uuid).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("%w: %w", ErrOrderNotFound, err)
	} else if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrGettingOrder, err)
	}

	return &order, nil
}

func (*RepositoryImpl) GetActiveOrderByMenu(tx *gorm.DB, menuUuid *uuid.UUID) (*Order, error) {
	var order Order
	err := tx.Model(&Order{}).Preload("Items").Where(&Order{MenuUuid: menuUuid}).First(&order).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("%w: %w", ErrOrderNotFound, err)
	} else if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrGettingOrder, err)
	}

	return &order, nil
}

func (*RepositoryImpl) GetActiveOrderByMenuName(tx *gorm.DB, menuName string) (*Order, error) {
	var order Order
	err := tx.Model(&Order{}).Preload("Items").Joins("JOIN menus ON menus.uuid = orders.menu_uuid").Where("menus.name = ?", menuName).First(&order).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("%w: %w", ErrOrderNotFound, err)
	} else if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrGettingOrder, err)
	}

	return &order, nil
}

func (*RepositoryImpl) GetAllOrderItems(tx *gorm.DB, orderUuid *uuid.UUID) ([]OrderItem, error) {
	orderItems := []OrderItem{}
	err := tx.Find(&orderItems).Error
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrCannotGetAllOrders, err)
	}

	return orderItems, nil
}

func (*RepositoryImpl) GetOrderItem(tx *gorm.DB, uuid *uuid.UUID) (*OrderItem, error) {
	orderItem := OrderItem{}
	err := tx.First(&orderItem, uuid).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("%w: %w", ErrOrderItemNotFound, err)
	} else if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrGettingOrderItem, err)
	}

	return &orderItem, nil
}

func (r *RepositoryImpl) CreateOrderItem(tx *gorm.DB, order_uuid *uuid.UUID, orderItem *OrderItem) (*OrderItem, error) {
	menuItemUuid := orderItem.MenuItemUuid
	if menuItemUuid == nil {
		return nil, fmt.Errorf("%w: %w", ErrCreatingOrderItem, ErrMenuItemUuidMissing)
	}

	menuItem, err := r.GetMenuItem(tx, menuItemUuid)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrCreatingOrderItem, err)
	}

	orderItem.Paid = false
	orderItem.Price = menuItem.Price

	err = tx.Create(&orderItem).Error
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrCreatingOrderItem, err)
	}

	return orderItem, nil
}

func (repo *RepositoryImpl) CreateOrder(tx *gorm.DB, order *Order) (*Order, error) {
	order.State = Open

	err := tx.Create(&order).Error
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrCreatingOrder, err)
	}

	return order, nil
}

func (repo *RepositoryImpl) UpdateOrder(tx *gorm.DB, orderUuid *uuid.UUID, order *Order) (*Order, error) {
	existingOrder, err := repo.GetOrder(tx, orderUuid)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrUpdatingOrder, err)
	}

	switch existingOrder.State {
	case Open:
		existingOrder.OrderDeadline = order.OrderDeadline
		if order.State == Finalized {
			existingOrder.State = Finalized
		} else if order.State != existingOrder.State {
			return nil, fmt.Errorf("%w: %w: from %s to %s", ErrUpdatingOrder, ErrOrderStateTransitionInvalid, existingOrder.State, order.State)
		}
		break
	case Finalized:
		if order.State == Ordered {
			existingOrder.State = Ordered
		} else if order.State != existingOrder.State {
			return nil, fmt.Errorf("%w: %w: from %s to %s", ErrUpdatingOrder, ErrOrderStateTransitionInvalid, existingOrder.State, order.State)
		}
		break
	case Ordered:
		existingOrder.Eta = order.Eta
		if order.State == Delivered {
			existingOrder.State = Delivered
		} else if order.State != existingOrder.State {
			return nil, fmt.Errorf("%w: %w: from %s to %s", ErrUpdatingOrder, ErrOrderStateTransitionInvalid, existingOrder.State, order.State)
		}
		break
	}

	existingOrder.SugarPerson = order.SugarPerson
	//if admin {
	//	existingOrder.State = order.State
	//}

	err = tx.Save(existingOrder).Error
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrUpdatingOrder, err)
	}

	return existingOrder, nil
}

func (repo *RepositoryImpl) UpdateOrderItem(tx *gorm.DB, orderItemUuid *uuid.UUID, orderItem *OrderItem) (*OrderItem, error) {
	existingOrderItem, err := repo.GetOrderItem(tx, orderItemUuid)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrUpdatingOrderItem, err)
	}

	existingOrderItem.User = orderItem.User
	existingOrderItem.Price = orderItem.Price
	existingOrderItem.Paid = orderItem.Paid
	existingOrderItem.MenuItemUuid = orderItem.MenuItemUuid

	err = tx.Save(existingOrderItem).Error
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrUpdatingOrderItem, err)
	}

	return existingOrderItem, nil
}

func (repo *RepositoryImpl) DeleteOrderItem(tx *gorm.DB, orderItemUuid *uuid.UUID) error {
	err := tx.Delete(&OrderItem{}, orderItemUuid).Error
	if err != nil {
		return fmt.Errorf("%w: %w", ErrDeletingOrderItem, err)
	}

	return nil
}

func (repo *RepositoryImpl) DeleteOrder(tx *gorm.DB, orderUuid *uuid.UUID) error {
	err := tx.Delete(&Order{}, orderUuid).Error
	if err != nil {
		return fmt.Errorf("%: %w", ErrDeletingOrderItem, err)
	}

	return nil
}
