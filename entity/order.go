package entity

import (
	"fmt"
	"time"

	"github.com/gofrs/uuid"
	"github.com/jmoiron/sqlx"
)

type Order struct {
	Uuid          uuid.UUID `db:"uuid"`
	Initiator     uuid.UUID `db:"initator"`
	SugarPerson   uuid.UUID `db:"sugar_person"`
	State         string    `db:"state"`
	OrderDeadline time.Time `db:"order_deadline"`
	Eta           time.Time `db:"eta"`
	MenuUuid      uuid.UUID `db:"menu_uuid"`
}

type OrderWithItems struct {
	Uuid          uuid.UUID `db:"uuid"`
	Initiator     uuid.UUID `db:"initator"`
	SugarPerson   uuid.UUID `db:"sugar_person"`
	State         string    `db:"state"`
	OrderDeadline time.Time `db:"order_deadline"`
	Eta           time.Time `db:"eta"`
	MenuUuid      uuid.UUID `db:"menu_uuid"`
	Items         []OrderItem
}

type OrderItem struct {
	Uuid         uuid.UUID `db:"uuid"`
	Price        int       `db:"price"`
	User         uuid.UUID `db:"order_user"`
	OrderUuid    uuid.UUID `db:"order_uuid"`
	MenuItemUuid uuid.UUID `db:"menu_item_uuid"`
}

type NewOrder struct {
	Initiator     uuid.UUID
	SugarPerson   uuid.UUID
	State         string
	OrderDeadline time.Time
	Eta           time.Time
	MenuUuid      uuid.UUID
}

type NewOrderItem struct {
	Price        int
	User         uuid.UUID
	OrderUuid    uuid.UUID
	MenuItemUuid uuid.UUID
}

func (*Repository) GetAllOrders(tx *sqlx.Tx) ([]OrderWithItems, error) {
	orders_map := map[uuid.UUID]*OrderWithItems{}
	rows, err := tx.Queryx("SELECT * FROM orders")
	if err != nil {
		return nil, fmt.Errorf("could not get all orders from db: %w", err)
	}
	for rows.Next() {
		var order OrderWithItems
		rows.StructScan(&order)
		orders_map[order.Uuid] = &order
	}

	rows, err = tx.Queryx("SELECT oi.* FROM orders o JOIN order_items oi on o.uuid = oi.order_uuid")
	if err != nil {
		return nil, fmt.Errorf("could not get all order_items from db: %w", err)
	}
	for rows.Next() {
		var order_item OrderItem
		rows.StructScan(&order_item)
		orders_map[order_item.OrderUuid].Items = append(orders_map[order_item.OrderUuid].Items, order_item)
	}

	orders := make([]OrderWithItems, 0, len(orders_map))
	for _, value := range orders_map {
		orders = append(orders, *value)
	}

	return orders, nil
}

func (*Repository) GetOrderWithItems(tx sqlx.Tx, uuid uuid.UUID) (*OrderWithItems, error) {
	var order OrderWithItems
	if err := tx.Get(&order, "SELECT * FROM orders WHERE uuid=?", uuid); err != nil {
		return nil, fmt.Errorf("error getting order %s: %w", uuid, err)
	}

	var orderItems []OrderItem
	if err := tx.Select(&orderItems, "SELECT * FROM order_items WHERE order_uuid=?", order.Uuid); err != nil {
		return nil, fmt.Errorf("error getting order items for order %s: %w", uuid, err)
	}

	order.Items = orderItems
	return &order, nil
}

func (*Repository) GetOrderItem(tx *sqlx.Tx, uuid uuid.UUID) (*OrderItem, error) {
	var order_item OrderItem
	if err := tx.Get(&order_item, "SELECT * FROM order_items WHERE id=?", uuid); err != nil {
		return nil, fmt.Errorf("error getting order item %s: %w", uuid, err)
	}

	return &order_item, nil
}

func (*Repository) CreateOrderItem(tx *sqlx.Tx, order_item NewOrderItem) (*OrderItem, error) {
	var uuid_string string
	err := tx.Get(
		&uuid_string,
		"INSERT INTO order_items (price, user_uuid, order_uuid, menu_item_uuid) VALUES ($1, $2, $3, $4) RETURNING uuid",
		order_item.Price, order_item.User, order_item.OrderUuid, order_item.MenuItemUuid,
	)
	if err != nil {
		return nil, fmt.Errorf("could not create order item: %w", err)
	}

	uuid := uuid.Must(uuid.FromString(uuid_string))

	return &OrderItem{
		Uuid:         uuid,
		Price:        order_item.Price,
		User:         order_item.User,
		OrderUuid:    order_item.OrderUuid,
		MenuItemUuid: order_item.MenuItemUuid,
	}, nil
}
