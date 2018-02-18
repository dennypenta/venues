package validator

import (
	"github.com/labstack/echo"
	"gopkg.in/go-playground/validator.v9"
)

// function for synchronize models.Restaurant.City
// we don't have to have cities like "Moscow" and "Mascow"
func ValidateCity(fl validator.FieldLevel) bool {
	// fake implementation
	// actually we have to have a separate service as Redis
	// that has a data set of cities
	if fl.Field().String() == "Mascow" {
		return false
	}

	return true
}

type Validator struct {
	validator *validator.Validate
}

func (v *Validator) Validate(i interface{}) error {
	return v.validator.Struct(i)
}

func NewValidator() echo.Validator {
	validatorType := validator.New()
	validatorType.RegisterValidation("city", ValidateCity)
	return &Validator{validator: validatorType}
}
