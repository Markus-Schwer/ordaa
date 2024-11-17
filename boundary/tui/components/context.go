package components

import (
	"context"
	"errors"

	"github.com/gofrs/uuid"
)

const (
	menuKey = "component_context_menu_uuid"
	USER_KEY
)

var (
	ErrContextValueNotSet = errors.New("context value not set")
	ErrContextValueUnexpectedType = errors.New("context value has unexpected type")
)

func SetMenuUuidToContext(ctx context.Context, orderId *uuid.UUID) context.Context {
	return context.WithValue(ctx, menuKey, *orderId)
}

func GetMenuUuidFromContext(ctx context.Context) (*uuid.UUID, error) {
	rawOrderId := ctx.Value(menuKey)
	if rawOrderId == nil {
		return nil, ErrContextValueNotSet
	}
	orderId, ok := rawOrderId.(uuid.UUID)
	if !ok {
		return nil, ErrContextValueUnexpectedType
	}
	return &orderId, nil
}
