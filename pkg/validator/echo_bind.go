package validator

import (
	"github.com/labstack/echo"
	"gopkg.in/go-playground/validator.v9"
)
func customFunc(fl FieldLevel) bool {

	if fl.Field().String() == "invalid" {
		return false
	}

	return true
}

validate.RegisterValidation("custom tag name", customFunc)


type Validator struct {
	validator *validator.Validate
}

func (v *Validator) Validate(i interface{}) error {
	return v.validator.Struct(i)
}

func NewValidator() echo.Validator {
	return &Validator{validator: validator.New()}
}
