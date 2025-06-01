package handler

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"maunium.net/go/mautrix/event"
	"maunium.net/go/mautrix/id"

	"github.com/Markus-Schwer/ordaa/internal/entity"
	"github.com/Markus-Schwer/ordaa/internal/repository"
)

func TestRegister(t *testing.T) {
	ctx := t.Context()

	type testCase struct {
		name        string
		sender      string
		msg         string
		userService UserService
		matches     bool
		response    *CommandResponse
	}

	testCases := []testCase{
		{
			name:   "should handle register command",
			sender: "@test:matrix.org",
			msg:    fmt.Sprintf("%s register", MatrixCommandPrefix),
			userService: &UserServiceMock{
				RegisterMatrixUserFunc: func(ctx context.Context, username string) (*entity.User, error) {
					return &entity.User{Name: username}, nil
				},
			},
			matches:  true,
			response: &CommandResponse{Msg: "successfully registered user: @test:matrix.org"},
		},
		{
			name:   "should handle register command error",
			sender: "@test:matrix.org",
			msg:    fmt.Sprintf("%s register", MatrixCommandPrefix),
			userService: &UserServiceMock{
				RegisterMatrixUserFunc: func(ctx context.Context, username string) (*entity.User, error) {
					return nil, repository.ErrCreatingUser
				},
			},
			matches:  true,
			response: &CommandResponse{Msg: "could not register user: could not create user"},
		},
		{
			name:    "should not match register command without prefix",
			msg:     "register",
			matches: false,
		},
		{
			name:    "should not match register command from empty message",
			msg:     MatrixCommandPrefix,
			matches: false,
		},
		{
			name:    "should not match register command with trailing whitespaces",
			msg:     fmt.Sprintf("%s register ", MatrixCommandPrefix),
			matches: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			h := RegisterHandler{
				UserService: tc.userService,
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
