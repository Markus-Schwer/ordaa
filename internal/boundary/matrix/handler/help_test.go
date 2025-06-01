package handler

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"maunium.net/go/mautrix/event"
)

func TestHelp(t *testing.T) {
	ctx := t.Context()
	h := HelpHandler{}

	type testCase struct {
		name     string
		msg      string
		matches  bool
		response *CommandResponse
	}

	testCases := []testCase{
		{
			name:     "should handle help command",
			msg:      fmt.Sprintf("%s help", MatrixCommandPrefix),
			matches:  true,
			response: &CommandResponse{Msg: "Hello world"},
		},
		{
			name:    "should not match help command without prefix",
			msg:     "help",
			matches: false,
		},
		{
			name:    "should not match help command from empty message",
			msg:     MatrixCommandPrefix,
			matches: false,
		},
		{
			name:    "should not match help command with trailing whitespaces",
			msg:     fmt.Sprintf("%s help ", MatrixCommandPrefix),
			matches: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			evt := &event.Event{
				Content: event.Content{
					Parsed: &event.MessageEventContent{
						Body: tc.msg,
					},
				},
			}

			matches := h.Matches(ctx, evt)
			assert.Equal(t, tc.matches, matches)

			if matches {
				resp := h.Handle(ctx, evt)

				if tc.response != nil {
					assert.NotNil(t, resp)
					assert.Equal(t, tc.response, resp)
				} else {
					assert.Nil(t, resp)
				}
			}
		})
	}
}
