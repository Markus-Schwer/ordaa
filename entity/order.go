package entity

import (
	"fmt"
	"time"

	"github.com/gofrs/uuid"
	"github.com/jmoiron/sqlx"
)

type Order struct {
	Uuid          uuid.UUID  `db:"uuid" json:"uuid"`
	Initiator     uuid.UUID  `db:"initiator" json:"initiator"`
	SugarPerson   *uuid.UUID `db:"sugar_person" json:"sugar_person"`
	State         string     `db:"state" json:"state"`
	OrderDeadline *time.Time `db:"order_deadline" json:"order_deadline"`
	Eta           *time.Time `db:"eta" json:"eta"`
	MenuUuid      uuid.UUID  `db:"menu_uuid" json:"menu_uuid"`
}

type OrderWithItems struct {
	Uuid          uuid.UUID   `db:"uuid" json:"uuid"`
	Initiator     uuid.UUID   `db:"initiator" json:"initiator"`
	SugarPerson   *uuid.UUID  `db:"sugar_person" json:"sugar_person"`
	State         string      `db:"state" json:"state"`
	OrderDeadline *time.Time  `db:"order_deadline" json:"order_deadline"`
	Eta           *time.Time  `db:"eta" json:"eta"`
	MenuUuid      uuid.UUID   `db:"menu_uuid" json:"menu_uuid"`
	Items         []OrderItem `db:"items" json:"items"`
}

type OrderItem struct {
	Uuid         uuid.UUID `db:"uuid" json:"uuid"`
	Price        int       `db:"price" json:"price"`
	Paid         bool      `db:"paid" json:"paid"`
	User         uuid.UUID `db:"order_user" json:"order_user"`
	OrderUuid    uuid.UUID `db:"order_uuid" json:"order_uuid"`
	MenuItemUuid uuid.UUID `db:"menu_item_uuid" json:"menu_item_uuid"`
}

type NewOrder struct {
	Initiator     uuid.UUID  `json:"initiator"`
	SugarPerson   *uuid.UUID `json:"sugar_person"`
	State         string     `json:"state"`
	OrderDeadline *time.Time `json:"order_deadline"`
	Eta           *time.Time `json:"eta"`
	MenuUuid      uuid.UUID  `json:"menu_uuid"`
}

type NewOrderItem struct {
	Price        int       `json:"price"`
	User         uuid.UUID `json:"user"`
	MenuItemUuid uuid.UUID `json:"menu_item_uuid"`
}

func (*Repository) GetAllOrders(tx *sqlx.Tx) ([]OrderWithItems, error) {
	ordersMap := map[uuid.UUID]*OrderWithItems{}
	rows, err := tx.Queryx("SELECT * FROM orders")
	if err != nil {
		return nil, fmt.Errorf("could not get all orders from db: %w", err)
	}
	for rows.Next() {
		var order OrderWithItems
		rows.StructScan(&order)
		ordersMap[order.Uuid] = &order
	}

	rows, err = tx.Queryx("SELECT oi.* FROM orders o JOIN order_items oi on o.uuid = oi.order_uuid")
	if err != nil {
		return nil, fmt.Errorf("could not get all order_items from db: %w", err)
	}
	for rows.Next() {
		var orderItem OrderItem
		rows.StructScan(&orderItem)
		ordersMap[orderItem.OrderUuid].Items = append(ordersMap[orderItem.OrderUuid].Items, orderItem)
	}

	orders := make([]OrderWithItems, 0, len(ordersMap))
	for _, value := range ordersMap {
		orders = append(orders, *value)
	}

	return orders, nil
}

func (*Repository) GetOrder(tx *sqlx.Tx, uuid uuid.UUID) (*OrderWithItems, error) {
	var order OrderWithItems
	if err := tx.Get(&order, "SELECT * FROM orders WHERE uuid=$1", uuid); err != nil {
		return nil, fmt.Errorf("error getting order %s: %w", uuid, err)
	}

	var orderItems []OrderItem
	if err := tx.Select(&orderItems, "SELECT * FROM order_items WHERE order_uuid=$1", order.Uuid); err != nil {
		return nil, fmt.Errorf("error getting order items for order %s: %w", uuid, err)
	}

	order.Items = orderItems
	return &order, nil
}

func (*Repository) GetOrderItem(tx *sqlx.Tx, uuid uuid.UUID) (*OrderItem, error) {
	var orderItem OrderItem
	if err := tx.Get(&orderItem, "SELECT * FROM order_items WHERE uuid=$1", uuid); err != nil {
		return nil, fmt.Errorf("error getting order item %s: %w", uuid, err)
	}

	return &orderItem, nil
}

func (*Repository) CreateOrderItem(tx *sqlx.Tx, order_uuid uuid.UUID, orderItem NewOrderItem) (*OrderItem, error) {
	var uuidString string
	err := tx.Get(
		&uuidString,
		"INSERT INTO order_items (price, user_uuid, menu_item_uuid, order_uuid) VALUES ($1, $2, $3, $4) RETURNING uuid",
		orderItem.Price, orderItem.User, orderItem.MenuItemUuid, order_uuid,
	)
	if err != nil {
		return nil, fmt.Errorf("could not create order item: %w", err)
	}

	uuid := uuid.Must(uuid.FromString(uuidString))

	return &OrderItem{
		Uuid:         uuid,
		Price:        orderItem.Price,
		User:         orderItem.User,
		OrderUuid:    order_uuid,
		MenuItemUuid: orderItem.MenuItemUuid,
	}, nil
}

func (repo *Repository) CreateOrder(tx *sqlx.Tx, order *NewOrder) (*Order, error) {
	var createdOrder Order
	err := tx.Get(
		&createdOrder,
		"INSERT INTO orders (initiator, sugar_person, state, order_deadline, eta, menu_uuid) VALUES ($1, $2, $3, $4, $5, $6) RETURNING *",
		order.Initiator, order.SugarPerson, order.State, order.OrderDeadline, order.Eta, order.MenuUuid,
	)
	if err != nil {
		return nil, fmt.Errorf("could not create order: %w", err)
	}

	return &createdOrder, nil
}

func (repo *Repository) UpdateOrder(tx *sqlx.Tx, orderUuid uuid.UUID, order *NewOrder) (*OrderWithItems, error) {
	_, err := tx.Exec(
		"UPDATE orders SET initiator = $2, sugar_person = $3, state = $4, order_deadline = $5, eta = $6, menu_uuid = $7 WHERE uuid = $1",
		orderUuid, order.Initiator, order.SugarPerson, order.State, order.OrderDeadline, order.Eta, order.MenuUuid,
	)
	if err != nil {
		return nil, fmt.Errorf("could not update order %s: %w", orderUuid, err)
	}

	return repo.GetOrder(tx, orderUuid)
}

func (repo *Repository) DeleteOrderItem(tx *sqlx.Tx, orderItemUuid uuid.UUID) error {
	_, err := tx.Exec("DELETE FROM order_items WHERE uuid = $1", orderItemUuid)
	if err != nil {
		return fmt.Errorf("could not delete order item %s: %w", orderItemUuid, err)
	}

	return nil
}

func (repo *Repository) DeleteOrder(tx *sqlx.Tx, orderUuid uuid.UUID) error {
	_, err := tx.Exec("DELETE FROM orders WHERE uuid = $1", orderUuid)
	if err != nil {
		return fmt.Errorf("could not delete order %s: %w", orderUuid, err)
	}

	return nil
}
