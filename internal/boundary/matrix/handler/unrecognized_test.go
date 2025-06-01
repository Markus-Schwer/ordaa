package handler

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"maunium.net/go/mautrix/event"
)

func TestUnrecognized(t *testing.T) {
	ctx := t.Context()
	h := UnrecognizedCommandHandler{}

	type testCase struct {
		name     string
		msg      string
		matches  bool
		response *CommandResponse
	}

	testCases := []testCase{
		{
			name:     "should handle empty command with prefix",
			msg:      MatrixCommandPrefix,
			matches:  true,
			response: &CommandResponse{Msg: fmt.Sprintf("command not recognized: %s", MatrixCommandPrefix)},
		},
		{
			name:     "should handle any command with prefix",
			msg:      fmt.Sprintf("%s asdf sadfk;", MatrixCommandPrefix),
			matches:  true,
			response: &CommandResponse{Msg: fmt.Sprintf("command not recognized: %s asdf sadfk;", MatrixCommandPrefix)},
		},
		{
			name:    "should not match without prefix",
			msg:     "asdf",
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
