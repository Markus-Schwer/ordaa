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

func TestAdd(t *testing.T) {
	ctx := t.Context()

	userUUID := uuid.Must(uuid.NewV4())

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
			name:   "should handle add command",
			sender: "@test:matrix.org",
			msg:    fmt.Sprintf("%s add sangam 62", MatrixCommandPrefix),
			userService: &UserServiceMock{
				GetMatrixUserByUsernameFunc: func(ctx context.Context, username string) (*entity.MatrixUser, error) {
					if username == "@test:matrix.org" {
						return &entity.MatrixUser{UserUUID: &userUUID}, nil
					}

					return nil, repository.ErrUserNotFound
				},
			},
			orderService: &OrderServiceMock{
				AddOrderItemToOrderByNameFunc: func(ctx context.Context, currentUser *uuid.UUID, shortName, menuName string) error {
					if menuName != "sangam" {
						return repository.ErrMenuNotFound
					}

					if shortName != "62" {
						return repository.ErrMenuItemNotFound
					}

					return nil
				},
			},
			matches:  true,
			response: &CommandResponse{Msg: "added 62 to active order sangam"},
		},
		{
			name:   "should handle add command user not found error",
			sender: "@test:matrix.org",
			msg:    fmt.Sprintf("%s add sangam 62", MatrixCommandPrefix),
			userService: &UserServiceMock{
				GetMatrixUserByUsernameFunc: func(ctx context.Context, username string) (*entity.MatrixUser, error) {
					return nil, repository.ErrUserNotFound
				},
			},
			matches:  true,
			response: &CommandResponse{Msg: "could not add to order: user not found"},
		},
		{
			name:    "should not match add command without prefix",
			msg:     "add",
			matches: false,
		},
		{
			name:    "should not match add command from empty message",
			msg:     MatrixCommandPrefix,
			matches: false,
		},
		{
			name:    "should not match add command without menu name",
			msg:     fmt.Sprintf("%s add ", MatrixCommandPrefix),
			matches: false,
		},
		{
			name:    "should not match add command without valid menu name",
			msg:     fmt.Sprintf("%s add sangam asdf 12345", MatrixCommandPrefix),
			matches: false,
		},
		{
			name:    "should not match add command without valid menu name",
			msg:     fmt.Sprintf("%s add san-gam 62", MatrixCommandPrefix),
			matches: false,
		},
		{
			name:    "should not match add command without valid menu name",
			msg:     fmt.Sprintf("%s add sangam 62-", MatrixCommandPrefix),
			matches: false,
		},
		{
			name:    "should not match add command with trailing whitespaces",
			msg:     fmt.Sprintf("%s add sangam ", MatrixCommandPrefix),
			matches: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			h := AddHandler{
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
