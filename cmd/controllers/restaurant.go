package controllers

import (
	"net/http"
	"venues/cmd/repositories"

	"venues/cmd/models"

	"github.com/labstack/echo"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
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

func (controller *RestaurantController) Create(context echo.Context) error {
	restaurant := &models.Restaurant{}
	if err := context.Bind(restaurant); err != nil {
		return context.String(http.StatusBadRequest, err.Error())
	}

	if err := controller.Repo.Create(restaurant); err != nil {
		return context.NoContent(http.StatusServiceUnavailable)
	}

	return context.NoContent(http.StatusOK)
}

func (controller *RestaurantController) Update(context echo.Context) error {
	query := &models.Restaurant{ID: bson.ObjectId(context.Param("restaurant_id"))}
	update := &models.Restaurant{}
	if err := context.Bind(update); err != nil {
		return context.String(http.StatusBadRequest, err.Error())
	}

	if err := controller.Repo.Update(query, update); err != nil {
		if err == mgo.ErrNotFound {
			return context.NoContent(http.StatusNotFound)
		}

		return context.NoContent(http.StatusServiceUnavailable)
	}

	return context.NoContent(http.StatusOK)

}
func NewRestaurantController() *RestaurantController {
	return &RestaurantController{Repo: repositories.NewRestaurantRepo()}
}
