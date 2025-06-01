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
	ErrCannotGetAllOrder               = errors.New("could not get all orders from db")
	ErrOrderNotFound                   = errors.New("order not found")
	ErrGettingOrder                    = errors.New("could not get order from db")
	ErrCannotGetAllOrders              = errors.New("could not get all orders from db")
	ErrCannotGetAllOrderItems          = errors.New("could not get all order items from db")
	ErrGettingOrderItemsOrderAndUser   = errors.New("could not get order items for order and user")
	ErrGettingOrderItem                = errors.New("error getting order item")
	ErrOrderItemNotFound               = errors.New("order item not found")
	ErrCreatingOrderItem               = errors.New("could not create order item")
	ErrOrderNotOpen                    = errors.New("order is not in state open")
	ErrCreatingOrder                   = errors.New("could not create order")
	ErrActiveOrderForMenuAlreadyExists = errors.New("there is already an active order the specified menu")
	ErrUpdatingOrder                   = errors.New("could not update order")
	ErrOrderStateTransitionInvalid     = errors.New("invalid order state transition")
	ErrUpdatingOrderItem               = errors.New("could not update order item")
	ErrOrderUUIDChangeForbidden        = errors.New("changing order uuid is forbidden")
	ErrDeletingOrderItem               = errors.New("could not delete order item")
	ErrPaidChangeForbidden             = errors.New("paid status can only be changed by sugar person")
	ErrMenuItemUUIDMissing             = errors.New("menu item uuid missing")
	ErrMenuItemUUIDChangeForbidden     = errors.New("changing menu item uuid is forbidden")
	ErrUserChangeForbidden             = errors.New("changing user is forbidden")
	ErrSugarPersonNotSet               = errors.New("the sugar person has not been set")
)

type OrderRepository struct {
	DB             *gorm.DB
	MenuRepository MenuRepository
}

func (r *OrderRepository) GetAllOrders(ctx context.Context) ([]entity.Order, error) {
	orders := []entity.Order{}

	err := r.DB.Model(&entity.Order{}).Find(&orders).Error
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrCannotGetAllOrder, err)
	}

	return orders, nil
}

func (r *OrderRepository) GetOrder(ctx context.Context, uuid *uuid.UUID) (*entity.Order, error) {
	var order entity.Order

	err := r.DB.Model(&entity.Order{}).First(&order, uuid).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("%w: %w", ErrOrderNotFound, err)
	} else if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrGettingOrder, err)
	}

	return &order, nil
}

func (r *OrderRepository) GetActiveOrderByMenu(ctx context.Context, menuUUID *uuid.UUID) (*entity.Order, error) {
	var order entity.Order

	err := r.DB.Model(&entity.Order{}).Where("menu_uuid = ? AND state != ?", menuUUID, entity.Delivered).First(&order).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("%w: %w", ErrOrderNotFound, err)
	} else if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrGettingOrder, err)
	}

	return &order, nil
}

func (r *OrderRepository) GetActiveOrderByMenuName(ctx context.Context, menuName string) (*entity.Order, error) {
	var order entity.Order

	err := r.DB.Model(&entity.Order{}).
		Joins("JOIN menus ON menus.uuid = orders.menu_uuid").
		Where("menus.name = ? AND state != ?", menuName, entity.Delivered).
		First(&order).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("%w: %w", ErrOrderNotFound, err)
	} else if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrGettingOrder, err)
	}

	return &order, nil
}

func (r *OrderRepository) GetAllOrderItems(ctx context.Context, orderUUID *uuid.UUID) ([]entity.OrderItem, error) {
	orderItems := []entity.OrderItem{}

	err := r.DB.Find(&orderItems).Error
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrCannotGetAllOrderItems, err)
	}

	return orderItems, nil
}

func (r *OrderRepository) GetAllOrderItemsForOrderAndUser(ctx context.Context, orderUUID, userUUID *uuid.UUID) ([]entity.OrderItem, error) {
	orderItems := []entity.OrderItem{}

	err := r.DB.Where(&entity.OrderItem{OrderUUID: orderUUID, User: userUUID}).Find(&orderItems).Error
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrGettingOrderItemsOrderAndUser, err)
	}

	return orderItems, nil
}

func (r *OrderRepository) GetOrderItem(ctx context.Context, uuid *uuid.UUID) (*entity.OrderItem, error) {
	orderItem := entity.OrderItem{}

	err := r.DB.First(&orderItem, uuid).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("%w: %w", ErrOrderItemNotFound, err)
	} else if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrGettingOrderItem, err)
	}

	return &orderItem, nil
}

func (r *OrderRepository) CreateOrderItem(
	ctx context.Context,
	orderUUID *uuid.UUID,
	orderItem *entity.OrderItem,
) (*entity.OrderItem, error) {
	tx := r.DB.Begin()

	menuItemUUID := orderItem.MenuItemUUID
	if menuItemUUID == nil {
		_ = tx.Rollback()
		return nil, fmt.Errorf("%w: %w", ErrCreatingOrderItem, ErrMenuItemUUIDMissing)
	}

	order, err := r.GetOrder(ctx, orderUUID)
	if err != nil {
		_ = tx.Rollback()
		return nil, fmt.Errorf("%w: %w", ErrCreatingOrderItem, err)
	}

	if order.State != entity.Open {
		_ = tx.Rollback()
		return nil, fmt.Errorf("%w: %w", ErrCreatingOrderItem, ErrOrderNotOpen)
	}

	menuItem, err := r.MenuRepository.GetMenuItem(ctx, menuItemUUID)
	if err != nil {
		_ = tx.Rollback()
		return nil, fmt.Errorf("%w: %w", ErrCreatingOrderItem, err)
	}

	orderItem.Paid = false
	orderItem.Price = menuItem.Price

	err = tx.Create(&orderItem).Error
	if err != nil {
		_ = tx.Rollback()
		return nil, fmt.Errorf("%w: %w", ErrCreatingOrderItem, err)
	}

	_ = tx.Commit()

	return orderItem, nil
}

func (r *OrderRepository) CreateOrder(ctx context.Context, order *entity.Order) (*entity.Order, error) {
	tx := r.DB.Begin()

	order.State = entity.Open

	_, err := r.GetActiveOrderByMenu(ctx, order.MenuUUID)
	if err == nil {
		_ = tx.Rollback()
		return nil, fmt.Errorf("%w: %w", ErrCreatingOrder, ErrActiveOrderForMenuAlreadyExists)
	} else if !errors.Is(err, ErrOrderNotFound) {
		_ = tx.Rollback()
		return nil, fmt.Errorf("%w: %w", ErrCreatingOrder, err)
	}

	err = tx.Create(&order).Error
	if err != nil {
		_ = tx.Rollback()
		return nil, fmt.Errorf("%w: %w", ErrCreatingOrder, err)
	}

	_ = tx.Commit()

	return order, nil
}

func (r *OrderRepository) UpdateOrder(ctx context.Context, orderUUID, currentUser *uuid.UUID, order *entity.Order) (*entity.Order, error) {
	tx := r.DB.Begin()

	order.UUID = orderUUID

	if err := tx.Save(order).Error; err != nil {
		return nil, fmt.Errorf("%w: %w", ErrUpdatingOrder, err)
	}

	return order, nil
}

func (r *OrderRepository) UpdateOrderItem(
	ctx context.Context,
	orderItemUUID,
	userUUID *uuid.UUID,
	orderItem *entity.OrderItem,
) (*entity.OrderItem, error) {
	tx := r.DB.Begin()

	existingOrderItem, err := r.GetOrderItem(ctx, orderItemUUID)
	if err != nil {
		_ = tx.Rollback()
		return nil, fmt.Errorf("%w: %w", ErrUpdatingOrderItem, err)
	}

	if *existingOrderItem.OrderUUID != *orderItem.OrderUUID {
		_ = tx.Rollback()
		return nil, fmt.Errorf("%w: %w", ErrUpdatingOrderItem, ErrOrderUUIDChangeForbidden)
	}

	if *existingOrderItem.MenuItemUUID != *orderItem.MenuItemUUID {
		_ = tx.Rollback()
		return nil, fmt.Errorf("%w: %w", ErrUpdatingOrderItem, ErrMenuItemUUIDChangeForbidden)
	}

	if *existingOrderItem.User != *orderItem.User {
		_ = tx.Rollback()
		return nil, fmt.Errorf("%w: %w", ErrUpdatingOrderItem, ErrUserChangeForbidden)
	}

	order, err := r.GetOrder(ctx, existingOrderItem.OrderUUID)
	if err != nil {
		_ = tx.Rollback()
		return nil, fmt.Errorf("%w: %w", ErrUpdatingOrderItem, err)
	}

	// check if sugar persion is nil
	if order.SugarPerson == nil {
		_ = tx.Rollback()
		return nil, ErrSugarPersonNotSet
	}

	if existingOrderItem.Paid != orderItem.Paid && *userUUID != *order.SugarPerson {
		_ = tx.Rollback()
		return nil, fmt.Errorf("%w: %w", ErrUpdatingOrderItem, ErrPaidChangeForbidden)
	}

	existingOrderItem.Paid = orderItem.Paid

	err = tx.Save(existingOrderItem).Error
	if err != nil {
		_ = tx.Rollback()
		return nil, fmt.Errorf("%w: %w", ErrUpdatingOrderItem, err)
	}

	_ = tx.Commit()

	return existingOrderItem, nil
}

func (r *OrderRepository) DeleteOrderItem(ctx context.Context, orderItemUUID *uuid.UUID) error {
	tx := r.DB.Begin()

	err := tx.Delete(&entity.OrderItem{}, orderItemUUID).Error
	if err != nil {
		_ = tx.Rollback()
		return fmt.Errorf("%w: %w", ErrDeletingOrderItem, err)
	}

	_ = tx.Commit()

	return nil
}

func (r *OrderRepository) DeleteOrder(ctx context.Context, orderUUID *uuid.UUID) error {
	tx := r.DB.Begin()

	err := tx.Delete(&entity.Order{}, orderUUID).Error
	if err != nil {
		_ = tx.Rollback()
		return fmt.Errorf("%w: %w", ErrDeletingOrderItem, err)
	}

	_ = tx.Commit()

	return nil
}
