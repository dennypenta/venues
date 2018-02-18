package controllers

import (
	"net/http"
	"venues/cmd/repositories"

	"github.com/labstack/echo"
)

type RestaurantController struct {
	Repo repositories.RestaurantAccessor
}

// I'm not checking for empty list cause We actually don't wanna see 204,
// Easier will get empty list and 200
func (controller *RestaurantController) List(context echo.Context) error {
	restaurants, err := controller.Repo.List()
	if err != nil {
		return context.NoContent(http.StatusServiceUnavailable)
	}

	return context.JSON(http.StatusOK, restaurants)
}

func NewRestaurantController() *RestaurantController {
	return &RestaurantController{Repo: repositories.NewRestaurantRepo()}
}
