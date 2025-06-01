package handler

import (
	"context"
	"fmt"
	"regexp"

	"github.com/gofrs/uuid"
	"maunium.net/go/mautrix/event"

	"github.com/Markus-Schwer/ordaa/internal/entity"
)

var registerRegex = regexp.MustCompile(fmt.Sprintf("^%s register$", MatrixCommandPrefixRegex))

//go:generate go tool moq -rm -out user_service_mock.go . UserService

type UserService interface {
	GetAllUsers(ctx context.Context) ([]entity.User, error)
	GetUser(ctx context.Context, uuid *uuid.UUID) (*entity.User, error)
	CreateUser(ctx context.Context, user *entity.User) (*entity.User, error)
	UpdateUser(ctx context.Context, uuid *uuid.UUID, user *entity.User) (*entity.User, error)
	DeleteUser(ctx context.Context, uuid *uuid.UUID) error
	RegisterMatrixUser(ctx context.Context, username string) (*entity.User, error)
	SetPublicKey(ctx context.Context, userUUID *uuid.UUID, publicKey string) error
	GetMatrixUserByUsername(ctx context.Context, username string) (*entity.MatrixUser, error)
}

type RegisterHandler struct {
	UserService UserService
}

func (h *RegisterHandler) Matches(ctx context.Context, evt *event.Event) bool {
	msg := evt.Content.AsMessage().Body
	return registerRegex.MatchString(msg)
}

func (h *RegisterHandler) Handle(ctx context.Context, evt *event.Event) *CommandResponse {
	username := evt.Sender.String()

	user, err := h.UserService.RegisterMatrixUser(ctx, username)
	if err != nil {
		return &CommandResponse{Msg: fmt.Sprintf("could not register user: %s", err)}
	}

	return &CommandResponse{Msg: fmt.Sprintf("successfully registered user: %s", user.Name)}
}
