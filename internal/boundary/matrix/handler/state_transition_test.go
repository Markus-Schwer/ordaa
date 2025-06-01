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
	"github.com/Markus-Schwer/ordaa/internal/service"
)

var (
	userUUID  = uuid.Must(uuid.NewV4())
	orderUUID = uuid.Must(uuid.NewV4())
)

type testCase struct {
	name         string
	sender       string
	msg          string
	orderService OrderService
	userService  UserService
	matches      bool
	response     *CommandResponse
}

func testStateTransition(tc *testCase, t *testing.T) {
	ctx := t.Context()

	t.Run(tc.name, func(t *testing.T) {
		h := StateTransitionHandler{
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

func TestFinalizeStateTransition(t *testing.T) {
	testCases := []testCase{
		{
			name:   "should handle finalize command",
			sender: "@test:matrix.org",
			msg:    fmt.Sprintf("%s finalize sangam", MatrixCommandPrefix),
			userService: &UserServiceMock{
				GetMatrixUserByUsernameFunc: func(ctx context.Context, username string) (*entity.MatrixUser, error) {
					if username == "@test:matrix.org" {
						return &entity.MatrixUser{UserUUID: &userUUID}, nil
					}

					return nil, repository.ErrUserNotFound
				},
			},
			orderService: &OrderServiceMock{
				GetActiveOrderByMenuNameFunc: func(ctx context.Context, name string) (*entity.Order, error) {
					return &entity.Order{UUID: &orderUUID, State: entity.Open}, nil
				},
				UpdateOrderFunc: func(ctx context.Context, currentUser, uuidMoqParam *uuid.UUID, order *entity.Order) (*entity.Order, error) {
					if *uuidMoqParam != orderUUID {
						return nil, repository.ErrOrderNotFound
					}

					if order.State != entity.Finalized {
						return nil, service.ErrOrderStateTransitionInvalid
					}

					return order, nil
				},
			},
			matches:  true,
			response: &CommandResponse{Msg: "successfully set state of order sangam to finalized"},
		},
		{
			name:   "should handle finalize command user not found error",
			sender: "@test:matrix.org",
			msg:    fmt.Sprintf("%s finalize sangam", MatrixCommandPrefix),
			userService: &UserServiceMock{
				GetMatrixUserByUsernameFunc: func(ctx context.Context, username string) (*entity.MatrixUser, error) {
					return nil, repository.ErrUserNotFound
				},
			},
			matches:  true,
			response: &CommandResponse{Msg: "could not update order: user not found"},
		},
		{
			name:   "should handle finalize command invalid state transition error",
			sender: "@test:matrix.org",
			msg:    fmt.Sprintf("%s finalize sangam", MatrixCommandPrefix),
			userService: &UserServiceMock{
				GetMatrixUserByUsernameFunc: func(ctx context.Context, username string) (*entity.MatrixUser, error) {
					return &entity.MatrixUser{UserUUID: &userUUID}, nil
				},
			},
			orderService: &OrderServiceMock{
				GetActiveOrderByMenuNameFunc: func(ctx context.Context, name string) (*entity.Order, error) {
					return &entity.Order{UUID: &orderUUID, State: entity.Open}, nil
				},
				UpdateOrderFunc: func(ctx context.Context, currentUser, uuidMoqParam *uuid.UUID, order *entity.Order) (*entity.Order, error) {
					return nil, service.ErrOrderStateTransitionInvalid
				},
			},
			matches:  true,
			response: &CommandResponse{Msg: "could not update order: invalid order state transition"},
		},
		{
			name:    "should not match finalize command without prefix",
			msg:     "finalize",
			matches: false,
		},
		{
			name:    "should not match finalize command from empty message",
			msg:     MatrixCommandPrefix,
			matches: false,
		},
		{
			name:    "should not match finalize command without menu name",
			msg:     fmt.Sprintf("%s finalize ", MatrixCommandPrefix),
			matches: false,
		},
		{
			name:    "should not match finalize command without valid menu name",
			msg:     fmt.Sprintf("%s finalize sangam 12345", MatrixCommandPrefix),
			matches: false,
		},
		{
			name:    "should not match finalize command without valid menu name",
			msg:     fmt.Sprintf("%s finalize san-gam", MatrixCommandPrefix),
			matches: false,
		},
		{
			name:    "should not match finalize command with trailing whitespaces",
			msg:     fmt.Sprintf("%s finalize sangam ", MatrixCommandPrefix),
			matches: false,
		},
	}

	for _, tc := range testCases {
		testStateTransition(&tc, t)
	}
}

func TestReOpenStateTransition(t *testing.T) {
	testCases := []testCase{
		{
			name:   "should handle re-open command",
			sender: "@test:matrix.org",
			msg:    fmt.Sprintf("%s re-open sangam", MatrixCommandPrefix),
			userService: &UserServiceMock{
				GetMatrixUserByUsernameFunc: func(ctx context.Context, username string) (*entity.MatrixUser, error) {
					if username == "@test:matrix.org" {
						return &entity.MatrixUser{UserUUID: &userUUID}, nil
					}

					return nil, repository.ErrUserNotFound
				},
			},
			orderService: &OrderServiceMock{
				GetActiveOrderByMenuNameFunc: func(ctx context.Context, name string) (*entity.Order, error) {
					return &entity.Order{UUID: &orderUUID, State: entity.Open}, nil
				},
				UpdateOrderFunc: func(ctx context.Context, currentUser, uuidMoqParam *uuid.UUID, order *entity.Order) (*entity.Order, error) {
					if *uuidMoqParam != orderUUID {
						return nil, repository.ErrUserNotFound
					}

					if order.State != entity.Open {
						return nil, service.ErrOrderStateTransitionInvalid
					}

					return order, nil
				},
			},
			matches:  true,
			response: &CommandResponse{Msg: "successfully set state of order sangam to open"},
		},
		{
			name:   "should handle re-open command user not found error",
			sender: "@test:matrix.org",
			msg:    fmt.Sprintf("%s re-open sangam", MatrixCommandPrefix),
			userService: &UserServiceMock{
				GetMatrixUserByUsernameFunc: func(ctx context.Context, username string) (*entity.MatrixUser, error) {
					return nil, repository.ErrUserNotFound
				},
			},
			matches:  true,
			response: &CommandResponse{Msg: "could not update order: user not found"},
		},
		{
			name:   "should handle re-open command invalid state transition error",
			sender: "@test:matrix.org",
			msg:    fmt.Sprintf("%s re-open sangam", MatrixCommandPrefix),
			userService: &UserServiceMock{
				GetMatrixUserByUsernameFunc: func(ctx context.Context, username string) (*entity.MatrixUser, error) {
					return &entity.MatrixUser{UserUUID: &userUUID}, nil
				},
			},
			orderService: &OrderServiceMock{
				GetActiveOrderByMenuNameFunc: func(ctx context.Context, name string) (*entity.Order, error) {
					return &entity.Order{UUID: &orderUUID, State: entity.Open}, nil
				},
				UpdateOrderFunc: func(ctx context.Context, currentUser, uuidMoqParam *uuid.UUID, order *entity.Order) (*entity.Order, error) {
					return nil, service.ErrOrderStateTransitionInvalid
				},
			},
			matches:  true,
			response: &CommandResponse{Msg: "could not update order: invalid order state transition"},
		},
		{
			name:    "should not match re-open command without prefix",
			msg:     "re-open",
			matches: false,
		},
		{
			name:    "should not match re-open command from empty message",
			msg:     MatrixCommandPrefix,
			matches: false,
		},
		{
			name:    "should not match re-open command without menu name",
			msg:     fmt.Sprintf("%s re-open ", MatrixCommandPrefix),
			matches: false,
		},
		{
			name:    "should not match re-open command without valid menu name",
			msg:     fmt.Sprintf("%s re-open sangam 12345", MatrixCommandPrefix),
			matches: false,
		},
		{
			name:    "should not match re-open command without valid menu name",
			msg:     fmt.Sprintf("%s re-open san-gam", MatrixCommandPrefix),
			matches: false,
		},
		{
			name:    "should not match re-open command with trailing whitespaces",
			msg:     fmt.Sprintf("%s re-open sangam ", MatrixCommandPrefix),
			matches: false,
		},
	}

	for _, tc := range testCases {
		testStateTransition(&tc, t)
	}
}

func TestOrderedStateTransition(t *testing.T) {
	testCases := []testCase{
		{
			name:   "should handle ordered command",
			sender: "@test:matrix.org",
			msg:    fmt.Sprintf("%s ordered sangam", MatrixCommandPrefix),
			userService: &UserServiceMock{
				GetMatrixUserByUsernameFunc: func(ctx context.Context, username string) (*entity.MatrixUser, error) {
					if username == "@test:matrix.org" {
						return &entity.MatrixUser{UserUUID: &userUUID}, nil
					}

					return nil, repository.ErrUserNotFound
				},
			},
			orderService: &OrderServiceMock{
				GetActiveOrderByMenuNameFunc: func(ctx context.Context, name string) (*entity.Order, error) {
					return &entity.Order{UUID: &orderUUID, State: entity.Open}, nil
				},
				UpdateOrderFunc: func(ctx context.Context, currentUser, uuidMoqParam *uuid.UUID, order *entity.Order) (*entity.Order, error) {
					if *uuidMoqParam != orderUUID {
						return nil, repository.ErrOrderNotFound
					}

					if order.State != entity.Ordered {
						return nil, service.ErrOrderStateTransitionInvalid
					}

					return order, nil
				},
			},
			matches:  true,
			response: &CommandResponse{Msg: "successfully set state of order sangam to ordered"},
		},
		{
			name:   "should handle ordered command user not found error",
			sender: "@test:matrix.org",
			msg:    fmt.Sprintf("%s ordered sangam", MatrixCommandPrefix),
			userService: &UserServiceMock{
				GetMatrixUserByUsernameFunc: func(ctx context.Context, username string) (*entity.MatrixUser, error) {
					return nil, repository.ErrUserNotFound
				},
			},
			matches:  true,
			response: &CommandResponse{Msg: "could not update order: user not found"},
		},
		{
			name:   "should handle ordered command invalid state transition error",
			sender: "@test:matrix.org",
			msg:    fmt.Sprintf("%s ordered sangam", MatrixCommandPrefix),
			userService: &UserServiceMock{
				GetMatrixUserByUsernameFunc: func(ctx context.Context, username string) (*entity.MatrixUser, error) {
					return &entity.MatrixUser{UserUUID: &userUUID}, nil
				},
			},
			orderService: &OrderServiceMock{
				GetActiveOrderByMenuNameFunc: func(ctx context.Context, name string) (*entity.Order, error) {
					return &entity.Order{UUID: &orderUUID, State: entity.Open}, nil
				},
				UpdateOrderFunc: func(ctx context.Context, currentUser, uuidMoqParam *uuid.UUID, order *entity.Order) (*entity.Order, error) {
					return nil, service.ErrOrderStateTransitionInvalid
				},
			},
			matches:  true,
			response: &CommandResponse{Msg: "could not update order: invalid order state transition"},
		},
		{
			name:    "should not match ordered command without prefix",
			msg:     "ordered",
			matches: false,
		},
		{
			name:    "should not match ordered command from empty message",
			msg:     MatrixCommandPrefix,
			matches: false,
		},
		{
			name:    "should not match ordered command without menu name",
			msg:     fmt.Sprintf("%s ordered ", MatrixCommandPrefix),
			matches: false,
		},
		{
			name:    "should not match ordered command without valid menu name",
			msg:     fmt.Sprintf("%s ordered sangam 12345", MatrixCommandPrefix),
			matches: false,
		},
		{
			name:    "should not match ordered command without valid menu name",
			msg:     fmt.Sprintf("%s ordered san-gam", MatrixCommandPrefix),
			matches: false,
		},
		{
			name:    "should not match ordered command with trailing whitespaces",
			msg:     fmt.Sprintf("%s ordered sangam ", MatrixCommandPrefix),
			matches: false,
		},
	}

	for _, tc := range testCases {
		testStateTransition(&tc, t)
	}
}

func TestDeliveredStateTransition(t *testing.T) {
	testCases := []testCase{
		{
			name:   "should handle delivered command",
			sender: "@test:matrix.org",
			msg:    fmt.Sprintf("%s delivered sangam", MatrixCommandPrefix),
			userService: &UserServiceMock{
				GetMatrixUserByUsernameFunc: func(ctx context.Context, username string) (*entity.MatrixUser, error) {
					if username == "@test:matrix.org" {
						return &entity.MatrixUser{UserUUID: &userUUID}, nil
					}

					return nil, repository.ErrUserNotFound
				},
			},
			orderService: &OrderServiceMock{
				GetActiveOrderByMenuNameFunc: func(ctx context.Context, name string) (*entity.Order, error) {
					return &entity.Order{UUID: &orderUUID, State: entity.Open}, nil
				},
				UpdateOrderFunc: func(ctx context.Context, currentUser, uuidMoqParam *uuid.UUID, order *entity.Order) (*entity.Order, error) {
					if *uuidMoqParam != orderUUID {
						return nil, repository.ErrOrderNotFound
					}

					if order.State != entity.Delivered {
						return nil, service.ErrOrderStateTransitionInvalid
					}

					return order, nil
				},
			},
			matches:  true,
			response: &CommandResponse{Msg: "successfully set state of order sangam to delivered"},
		},
		{
			name:   "should handle delivered command user not found error",
			sender: "@test:matrix.org",
			msg:    fmt.Sprintf("%s delivered sangam", MatrixCommandPrefix),
			userService: &UserServiceMock{
				GetMatrixUserByUsernameFunc: func(ctx context.Context, username string) (*entity.MatrixUser, error) {
					return nil, repository.ErrUserNotFound
				},
			},
			matches:  true,
			response: &CommandResponse{Msg: "could not update order: user not found"},
		},
		{
			name:   "should handle delivered command invalid state transition error",
			sender: "@test:matrix.org",
			msg:    fmt.Sprintf("%s delivered sangam", MatrixCommandPrefix),
			userService: &UserServiceMock{
				GetMatrixUserByUsernameFunc: func(ctx context.Context, username string) (*entity.MatrixUser, error) {
					return &entity.MatrixUser{UserUUID: &userUUID}, nil
				},
			},
			orderService: &OrderServiceMock{
				GetActiveOrderByMenuNameFunc: func(ctx context.Context, name string) (*entity.Order, error) {
					return &entity.Order{UUID: &orderUUID, State: entity.Open}, nil
				},
				UpdateOrderFunc: func(ctx context.Context, currentUser, uuidMoqParam *uuid.UUID, order *entity.Order) (*entity.Order, error) {
					return nil, service.ErrOrderStateTransitionInvalid
				},
			},
			matches:  true,
			response: &CommandResponse{Msg: "could not update order: invalid order state transition"},
		},
		{
			name:    "should not match delivered command without prefix",
			msg:     "delivered",
			matches: false,
		},
		{
			name:    "should not match delivered command from empty message",
			msg:     MatrixCommandPrefix,
			matches: false,
		},
		{
			name:    "should not match delivered command without menu name",
			msg:     fmt.Sprintf("%s delivered ", MatrixCommandPrefix),
			matches: false,
		},
		{
			name:    "should not match delivered command without valid menu name",
			msg:     fmt.Sprintf("%s delivered sangam 12345", MatrixCommandPrefix),
			matches: false,
		},
		{
			name:    "should not match delivered command without valid menu name",
			msg:     fmt.Sprintf("%s delivered san-gam", MatrixCommandPrefix),
			matches: false,
		},
		{
			name:    "should not match delivered command with trailing whitespaces",
			msg:     fmt.Sprintf("%s delivered sangam ", MatrixCommandPrefix),
			matches: false,
		},
	}

	for _, tc := range testCases {
		testStateTransition(&tc, t)
	}
}
