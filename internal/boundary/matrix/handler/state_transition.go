package handler

import (
	"context"
	"fmt"
	"regexp"

	"maunium.net/go/mautrix/event"

	"github.com/Markus-Schwer/ordaa/internal/entity"
)

var stateTransitionRegex = regexp.MustCompile(fmt.Sprintf("^%s (finalize|re-open|ordered|delivered) (\\w+)$", MatrixCommandPrefixRegex))

type StateTransitionHandler struct {
	UserService  UserService
	OrderService OrderService
}

func (h *StateTransitionHandler) Matches(ctx context.Context, evt *event.Event) bool {
	msg := evt.Content.AsMessage().Body

	return stateTransitionRegex.MatchString(msg)
}

func (h *StateTransitionHandler) Handle(ctx context.Context, evt *event.Event) *CommandResponse {
	currentUser, err := h.UserService.GetMatrixUserByUsername(ctx, evt.Sender.String())
	if err != nil {
		return &CommandResponse{Msg: fmt.Sprintf("could not update order: %s", err)}
	}

	msg := evt.Content.AsMessage().Body

	match := stateTransitionRegex.FindStringSubmatch(msg)
	if match == nil {
		return &CommandResponse{Msg: "could not update order: no menu name provided"}
	}

	action := match[1]
	menuName := match[2]

	order, err := h.OrderService.GetActiveOrderByMenuName(ctx, menuName)
	if err != nil {
		return &CommandResponse{Msg: fmt.Sprintf("could not update order: %s", err)}
	}

	switch action {
	case "finalize":
		order.State = entity.Finalized
	case "re-open":
		order.State = entity.Open
	case "ordered":
		order.State = entity.Ordered
	case "delivered":
		order.State = entity.Delivered
	}

	if _, err = h.OrderService.UpdateOrder(ctx, currentUser.UserUUID, order.UUID, order); err != nil {
		return &CommandResponse{Msg: fmt.Sprintf("could not update order: %s", err)}
	}

	return &CommandResponse{Msg: fmt.Sprintf("successfully set state of order %s to %s", menuName, order.State)}
}
