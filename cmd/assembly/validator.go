package assembly

import (
	"gopkg.in/go-playground/validator.v9"
	"github.com/labstack/echo"
)

type Validator struct {
	validator *validator.Validate
}

func (v *Validator) Validate(i interface{}) error {
	return v.validator.Struct(i)
}

func newValidator() echo.Validator {
	return &Validator{validator: validator.New()}
}

