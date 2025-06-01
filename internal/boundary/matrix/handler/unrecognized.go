package handler

import (
	"context"
	"fmt"
	"strings"

	"maunium.net/go/mautrix/event"
)

type UnrecognizedCommandHandler struct{}

func (h *UnrecognizedCommandHandler) Matches(ctx context.Context, evt *event.Event) bool {
	msg := evt.Content.AsMessage().Body
	return strings.HasPrefix(msg, MatrixCommandPrefix)
}

func (h *UnrecognizedCommandHandler) Handle(ctx context.Context, evt *event.Event) *CommandResponse {
	msg := evt.Content.AsMessage().Body
	return &CommandResponse{Msg: fmt.Sprintf("command not recognized: %s", msg)}
}
