package handler

import (
	"context"
	"fmt"
	"regexp"

	"github.com/gofrs/uuid"
	"maunium.net/go/mautrix/event"

	"github.com/Markus-Schwer/ordaa/internal/entity"
)

var startRegex = regexp.MustCompile(fmt.Sprintf("^%s start (\\w+)$", MatrixCommandPrefixRegex))

//go:generate go tool moq -rm -out order_service_mock.go . OrderService

type OrderService interface {
	GetAllOrders(ctx context.Context) ([]entity.Order, error)
	GetOrder(ctx context.Context, uuid *uuid.UUID) (*entity.Order, error)
	GetActiveOrderByMenu(ctx context.Context, menuUUID *uuid.UUID) (*entity.Order, error)
	GetActiveOrderByMenuName(ctx context.Context, name string) (*entity.Order, error)
	CreateOrder(ctx context.Context, currentUser *uuid.UUID, order *entity.Order) (*entity.Order, error)
	UpdateOrder(ctx context.Context, currentUser *uuid.UUID, uuid *uuid.UUID, order *entity.Order) (*entity.Order, error)
	CreateOrderForMenuName(ctx context.Context, currentUser *uuid.UUID, menuName string) (*entity.Order, error)
	AddOrderItemToOrderByName(ctx context.Context, currentUser *uuid.UUID, shortName, menuName string) error
}

type StartHandler struct {
	OrderService OrderService
	UserService  UserService
}

func (h *StartHandler) Matches(ctx context.Context, evt *event.Event) bool {
	msg := evt.Content.AsMessage().Body

	return startRegex.MatchString(msg)
}

func (h *StartHandler) Handle(ctx context.Context, evt *event.Event) *CommandResponse {
	currentUser, err := h.UserService.GetMatrixUserByUsername(ctx, evt.Sender.String())
	if err != nil {
		return &CommandResponse{Msg: fmt.Sprintf("could not start order: %s", err)}
	}

	msg := evt.Content.AsMessage().Body

	match := startRegex.FindStringSubmatch(msg)
	if match == nil {
		return &CommandResponse{Msg: "could not start order: no menu name provided"}
	}

	menuName := match[1]

	order, err := h.OrderService.CreateOrderForMenuName(ctx, currentUser.UserUUID, menuName)
	if err != nil {
		return &CommandResponse{Msg: fmt.Sprintf("could not start order: %s", err)}
	}

	return &CommandResponse{Msg: fmt.Sprintf("started new order for %s (id: %s)", menuName, order.UUID.String())}
}
