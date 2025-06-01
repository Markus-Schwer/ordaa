package handler

import (
	"context"
	"fmt"
	"regexp"

	"maunium.net/go/mautrix/event"
)

var addRegex = regexp.MustCompile(fmt.Sprintf("^%s add (\\w+) (\\w+)$", MatrixCommandPrefixRegex))

type AddHandler struct {
	OrderService OrderService
	UserService  UserService
}

func (h *AddHandler) Matches(ctx context.Context, evt *event.Event) bool {
	msg := evt.Content.AsMessage().Body

	return addRegex.MatchString(msg)
}

func (h *AddHandler) Handle(ctx context.Context, evt *event.Event) *CommandResponse {
	currentUser, err := h.UserService.GetMatrixUserByUsername(ctx, evt.Sender.String())
	if err != nil {
		return &CommandResponse{Msg: fmt.Sprintf("could not add to order: %s", err)}
	}

	msg := evt.Content.AsMessage().Body

	match := addRegex.FindStringSubmatch(msg)
	if match == nil {
		return &CommandResponse{Msg: "message must be in the format 'add [menu_name] [short_name]'"}
	}

	menuName := match[1]
	shortName := match[2]

	if err = h.OrderService.AddOrderItemToOrderByName(ctx, currentUser.UserUUID, shortName, menuName); err != nil {
		return &CommandResponse{Msg: fmt.Sprintf("could not add order: %s", err)}
	}

	return &CommandResponse{Msg: fmt.Sprintf("added %s to active order %s", shortName, menuName)}
}
