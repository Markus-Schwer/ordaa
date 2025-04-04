package utils

import (
	"fmt"

	"github.com/gofrs/uuid"
	"github.com/labstack/echo/v4"
)

func UuidParam(c echo.Context, name string) (*uuid.UUID, error) {
	uuidString := c.Param(name)
	uuid, err := uuid.FromString(uuidString)
	if err != nil {
		return nil, fmt.Errorf("could not parse param %s as uuid: %w", name, err)
	}
	return &uuid, nil
}

func BindAndValidate(ctx echo.Context, v interface{}) error {
	if err := ctx.Bind(v); err != nil {
		return err
	}
	if err := ctx.Validate(v); err != nil {
		return err
	}
	return nil
}
