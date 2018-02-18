package controllers

import (
	"net/http"
	"venues/cmd/repositories"

	"github.com/labstack/echo"
)

type RestaurantController struct {
	Repo repositories.Repo
}

// I'm not checking for empty list cause We actually don't wanna see 204,
// Easier will get empty list and 200
func (controller *RestaurantController) List(context echo.Context) error {
	restaurants := controller.Repo.List()
	return context.JSON(http.StatusOK, restaurants)
}
