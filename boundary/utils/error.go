package utils

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"net/http"
)

type Error struct {
	Errors map[string]interface{} `json:"errors"`
}

func WrapBindError(err error) error {
	switch v := err.(type) {
	case validator.ValidationErrors:
		msg := fmt.Sprintf("%s validation error: %s", v[0].Field(), v[0].Tag())
		return NewError(http.StatusUnprocessableEntity, msg)
	case *echo.HTTPError:
		return v
	default:
		return NewError(http.StatusInternalServerError, err.Error())
	}
}

func NewBindError(field, tag string) error {
	msg := fmt.Sprintf("%s validation error: %s", field, tag)
	return NewError(http.StatusUnprocessableEntity, msg)
}

func NewUnauthorized() error {
	return NewError(http.StatusUnauthorized, "Unauthorized")
}

func NewStatusUnprocessableEntity(msg string) error {
	return NewError(http.StatusUnprocessableEntity, msg)
}

func NewNotFoundError(msg string) error {
	if msg == "" {
		msg = "Not Found"
	}
	return NewError(http.StatusNotFound, msg)
}

func NewInternalServerError(err error) error {
	msg := "Interval Server Error"
	if err != nil {
		msg = err.Error()
	}
	return NewError(http.StatusInternalServerError, msg)
}

func NewError(statusCode int, msg string) error {
	return echo.NewHTTPError(statusCode, &Error{
		Errors: map[string]interface{}{
			"body": msg,
		},
	})
}
