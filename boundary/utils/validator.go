package utils

import "github.com/go-playground/validator/v10"

type Validator struct {
	delegate *validator.Validate
}

func NewValidator() *Validator {
	return &Validator{delegate: validator.New()}
}

func (v *Validator) Validate(i interface{}) error {
  return v.delegate.Struct(i)
}

