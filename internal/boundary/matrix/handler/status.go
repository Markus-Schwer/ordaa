package handler

import (
	"context"
	"fmt"
	"regexp"

	"maunium.net/go/mautrix/event"
)

var statusRegex = regexp.MustCompile(fmt.Sprintf("^%s status (\\w+)$", MatrixCommandPrefixRegex))

type StatusHandler struct {
	OrderService OrderService
}

func (h *StatusHandler) Matches(ctx context.Context, evt *event.Event) bool {
	msg := evt.Content.AsMessage().Body

	return statusRegex.MatchString(msg)
}

func (h *StatusHandler) Handle(ctx context.Context, evt *event.Event) *CommandResponse {
	msg := evt.Content.AsMessage().Body

	menuName := statusRegex.FindStringSubmatch(msg)[1]

	order, err := h.OrderService.GetActiveOrderByMenuName(ctx, menuName)
	if err != nil {
		return &CommandResponse{Msg: fmt.Sprintf("could not get status of order: %s", err)}
	}

	return &CommandResponse{Msg: order.State}
}
