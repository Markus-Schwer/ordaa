package handler

import (
	"context"
	"fmt"
	"testing"

	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/assert"
	"maunium.net/go/mautrix/event"
	"maunium.net/go/mautrix/id"

	"github.com/Markus-Schwer/ordaa/internal/entity"
	"github.com/Markus-Schwer/ordaa/internal/repository"
)

func TestStatus(t *testing.T) {
	ctx := t.Context()

	orderUUID := uuid.Must(uuid.NewV4())

	type testCase struct {
		name         string
		sender       string
		msg          string
		orderService OrderService
		matches      bool
		response     *CommandResponse
	}

	testCases := []testCase{
		{
			name:   "should handle status command",
			sender: "@test:matrix.org",
			msg:    fmt.Sprintf("%s status sangam", MatrixCommandPrefix),
			orderService: &OrderServiceMock{
				GetActiveOrderByMenuNameFunc: func(ctx context.Context, name string) (*entity.Order, error) {
					return &entity.Order{UUID: &orderUUID, State: entity.Open}, nil
				},
			},
			matches:  true,
			response: &CommandResponse{Msg: "open"},
		},
		{
			name:   "should handle status command order not found error",
			sender: "@test:matrix.org",
			msg:    fmt.Sprintf("%s status sangam", MatrixCommandPrefix),
			orderService: &OrderServiceMock{
				GetActiveOrderByMenuNameFunc: func(ctx context.Context, name string) (*entity.Order, error) {
					return nil, repository.ErrOrderNotFound
				},
			},
			matches:  true,
			response: &CommandResponse{Msg: "could not get status of order: order not found"},
		},
		{
			name:    "should not match status command without prefix",
			msg:     "status",
			matches: false,
		},
		{
			name:    "should not match status command from empty message",
			msg:     MatrixCommandPrefix,
			matches: false,
		},
		{
			name:    "should not match status command without menu name",
			msg:     fmt.Sprintf("%s status ", MatrixCommandPrefix),
			matches: false,
		},
		{
			name:    "should not match status command without valid menu name",
			msg:     fmt.Sprintf("%s status sangam 12345", MatrixCommandPrefix),
			matches: false,
		},
		{
			name:    "should not match status command without valid menu name",
			msg:     fmt.Sprintf("%s status san-gam", MatrixCommandPrefix),
			matches: false,
		},
		{
			name:    "should not match status command with trailing whitespaces",
			msg:     fmt.Sprintf("%s status sangam ", MatrixCommandPrefix),
			matches: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			h := StatusHandler{
				OrderService: tc.orderService,
			}

			evt := &event.Event{
				Sender: id.UserID(tc.sender),
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
