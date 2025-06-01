package handler

import (
	"context"
	"fmt"
	"regexp"

	"maunium.net/go/mautrix/event"
)

var helpRegex = regexp.MustCompile(fmt.Sprintf("^%s help$", MatrixCommandPrefixRegex))

type HelpHandler struct{}

func (h *HelpHandler) Matches(ctx context.Context, evt *event.Event) bool {
	msg := evt.Content.AsMessage().Body

	return helpRegex.MatchString(msg)
}

func (h *HelpHandler) Handle(ctx context.Context, evt *event.Event) *CommandResponse {
	return &CommandResponse{Msg: "Hello world"}
}
