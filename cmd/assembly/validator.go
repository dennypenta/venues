package assembly

import (
	"github.com/labstack/echo"
	"gopkg.in/go-playground/validator.v9"
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
