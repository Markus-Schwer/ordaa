package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/gofrs/uuid"

	"github.com/Markus-Schwer/ordaa/internal/entity"
	"github.com/Markus-Schwer/ordaa/internal/repository"
)

var (
	ErrCreatingOrder                   = errors.New("could not create order")
	ErrUpdatingOrder                   = errors.New("could not update order")
	ErrSugarPersonChangeForbidden      = errors.New("changing sugar person after it has already been set is forbidden")
	ErrOrderStateTransitionInvalid     = errors.New("invalid order state transition")
	ErrActiveOrderForMenuAlreadyExists = errors.New("there is already an active order the specified menu")
	ErrAddingOrderItem                 = errors.New("adding order item")
	ErrCannotReopenOrder               = errors.New("only the initiator can reopen the order")
)

type OrderRepository interface {
	GetAllOrders(ctx context.Context) ([]entity.Order, error)
	GetOrder(ctx context.Context, uuid *uuid.UUID) (*entity.Order, error)
	GetActiveOrderByMenu(ctx context.Context, menuUUID *uuid.UUID) (*entity.Order, error)
	GetActiveOrderByMenuName(ctx context.Context, menuName string) (*entity.Order, error)
	GetAllOrderItems(ctx context.Context, orderUUID *uuid.UUID) ([]entity.OrderItem, error)
	GetAllOrderItemsForOrderAndUser(ctx context.Context, orderUUID *uuid.UUID, userUUID *uuid.UUID) ([]entity.OrderItem, error)
	GetOrderItem(ctx context.Context, uuid *uuid.UUID) (*entity.OrderItem, error)
	CreateOrderItem(ctx context.Context, orderUUID *uuid.UUID, orderItem *entity.OrderItem) (*entity.OrderItem, error)
	CreateOrder(ctx context.Context, order *entity.Order) (*entity.Order, error)
	UpdateOrder(ctx context.Context, currentUser *uuid.UUID, orderUUID *uuid.UUID, order *entity.Order) (*entity.Order, error)
	UpdateOrderItem(ctx context.Context, orderItemUUID *uuid.UUID, userUUID *uuid.UUID, orderItem *entity.OrderItem) (*entity.OrderItem, error)
	DeleteOrderItem(ctx context.Context, orderItemUUID *uuid.UUID) error
	DeleteOrder(ctx context.Context, orderUUID *uuid.UUID) error
}

type OrderService struct {
	OrderRepository OrderRepository
	MenuRepository  MenuRepository
}

func (i *OrderService) GetAllOrders(ctx context.Context) ([]entity.Order, error) {
	return i.OrderRepository.GetAllOrders(ctx)
}

func (i *OrderService) GetOrder(ctx context.Context, uuid *uuid.UUID) (*entity.Order, error) {
	return i.OrderRepository.GetOrder(ctx, uuid)
}

func (i *OrderService) GetActiveOrderByMenu(ctx context.Context, menuUUID *uuid.UUID) (*entity.Order, error) {
	return i.OrderRepository.GetActiveOrderByMenu(ctx, menuUUID)
}

func (i *OrderService) GetActiveOrderByMenuName(ctx context.Context, name string) (*entity.Order, error) {
	return i.OrderRepository.GetActiveOrderByMenuName(ctx, name)
}

func (i *OrderService) CreateOrder(ctx context.Context, currentUser *uuid.UUID, order *entity.Order) (*entity.Order, error) {
	return i.OrderRepository.CreateOrder(ctx, order)
}

func (i *OrderService) UpdateOrder(ctx context.Context, currentUser, uuid *uuid.UUID, order *entity.Order) (*entity.Order, error) {
	existingOrder, err := i.GetOrder(ctx, uuid)
	if err != nil {
		return nil, err
	}

	switch existingOrder.State {
	case entity.Open:
		existingOrder.OrderDeadline = order.OrderDeadline
		if order.State == entity.Finalized {
			existingOrder.State = entity.Finalized
		} else if order.State != existingOrder.State {
			return nil, fmt.Errorf("%w: %w: from %s to %s", ErrUpdatingOrder, ErrOrderStateTransitionInvalid, existingOrder.State, order.State)
		}
	case entity.Finalized:
		if order.State == entity.Open && *currentUser != *existingOrder.Initiator {
			return nil, fmt.Errorf("%w: %w", ErrUpdatingOrder, ErrCannotReopenOrder)
		}

		if order.State == entity.Ordered || order.State == entity.Open {
			existingOrder.State = order.State
		} else if order.State != existingOrder.State {
			return nil, fmt.Errorf("%w: %w: from %s to %s", ErrUpdatingOrder, ErrOrderStateTransitionInvalid, existingOrder.State, order.State)
		}
	case entity.Ordered:
		existingOrder.Eta = order.Eta
		if order.State == entity.Delivered {
			existingOrder.State = entity.Delivered
		} else if order.State != existingOrder.State {
			return nil, fmt.Errorf("%w: %w: from %s to %s", ErrUpdatingOrder, ErrOrderStateTransitionInvalid, existingOrder.State, order.State)
		}
	}

	if existingOrder.SugarPerson != nil && *existingOrder.SugarPerson != *order.SugarPerson {
		return nil, fmt.Errorf("%w: %w", ErrUpdatingOrder, ErrSugarPersonChangeForbidden)
	}

	existingOrder.SugarPerson = order.SugarPerson

	return i.OrderRepository.UpdateOrder(ctx, currentUser, uuid, order)
}

func (i *OrderService) CreateOrderForMenuName(ctx context.Context, currentUser *uuid.UUID, menuName string) (*entity.Order, error) {
	menu, err := i.MenuRepository.GetMenuByName(ctx, menuName)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrCreatingOrder, err)
	}

	activeOrder, err := i.OrderRepository.GetActiveOrderByMenu(ctx, menu.UUID)
	if err != nil && !errors.Is(err, repository.ErrOrderNotFound) {
		return nil, fmt.Errorf("%w: %w", ErrCreatingOrder, err)
	}

	if activeOrder != nil {
		return nil, ErrActiveOrderForMenuAlreadyExists
	}

	order, err := i.OrderRepository.CreateOrder(ctx, &entity.Order{Initiator: currentUser, MenuUUID: menu.UUID})
	if err != nil {
		return nil, err
	}

	return order, nil
}

func (i *OrderService) AddOrderItemToOrderByName(ctx context.Context, currentUser *uuid.UUID, shortName, menuName string) error {
	order, err := i.OrderRepository.GetActiveOrderByMenuName(ctx, menuName)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrAddingOrderItem, err)
	}

	menuItem, err := i.MenuRepository.GetMenuItemByShortName(ctx, order.MenuUUID, shortName)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrAddingOrderItem, err)
	}

	orderItem := &entity.OrderItem{
		User:         currentUser,
		MenuItemUUID: menuItem.UUID,
		OrderUUID:    order.UUID,
	}

	if _, err = i.OrderRepository.CreateOrderItem(ctx, order.UUID, orderItem); err != nil {
		return fmt.Errorf("%w: %w", ErrAddingOrderItem, err)
	}

	return nil
}
