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

func TestStart(t *testing.T) {
	ctx := t.Context()

	userUUID := uuid.Must(uuid.NewV4())
	orderUUID := uuid.Must(uuid.NewV4())

	type testCase struct {
		name         string
		sender       string
		msg          string
		orderService OrderService
		userService  UserService
		matches      bool
		response     *CommandResponse
	}

	testCases := []testCase{
		{
			name:   "should handle start command",
			sender: "@test:matrix.org",
			msg:    fmt.Sprintf("%s start sangam", MatrixCommandPrefix),
			userService: &UserServiceMock{
				GetMatrixUserByUsernameFunc: func(ctx context.Context, username string) (*entity.MatrixUser, error) {
					if username == "@test:matrix.org" {
						return &entity.MatrixUser{UserUUID: &userUUID}, nil
					}

					return nil, repository.ErrUserNotFound
				},
			},
			orderService: &OrderServiceMock{
				CreateOrderForMenuNameFunc: func(ctx context.Context, currentUser *uuid.UUID, menuName string) (*entity.Order, error) {
					if menuName != "sangam" {
						return nil, repository.ErrMenuNotFound
					}

					return &entity.Order{UUID: &orderUUID}, nil
				},
			},
			matches:  true,
			response: &CommandResponse{Msg: fmt.Sprintf("started new order for sangam (id: %s)", orderUUID.String())},
		},
		{
			name:   "should handle start command user not found error",
			sender: "@test:matrix.org",
			msg:    fmt.Sprintf("%s start sangam", MatrixCommandPrefix),
			userService: &UserServiceMock{
				GetMatrixUserByUsernameFunc: func(ctx context.Context, username string) (*entity.MatrixUser, error) {
					return nil, repository.ErrUserNotFound
				},
			},
			matches:  true,
			response: &CommandResponse{Msg: "could not start order: user not found"},
		},
		{
			name:    "should not match start command without prefix",
			msg:     "start",
			matches: false,
		},
		{
			name:    "should not match start command from empty message",
			msg:     MatrixCommandPrefix,
			matches: false,
		},
		{
			name:    "should not match start command without menu name",
			msg:     fmt.Sprintf("%s start ", MatrixCommandPrefix),
			matches: false,
		},
		{
			name:    "should not match start command without valid menu name",
			msg:     fmt.Sprintf("%s start sangam 12345", MatrixCommandPrefix),
			matches: false,
		},
		{
			name:    "should not match start command without valid menu name",
			msg:     fmt.Sprintf("%s start san-gam", MatrixCommandPrefix),
			matches: false,
		},
		{
			name:    "should not match start command with trailing whitespaces",
			msg:     fmt.Sprintf("%s start sangam ", MatrixCommandPrefix),
			matches: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			h := StartHandler{
				UserService:  tc.userService,
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
