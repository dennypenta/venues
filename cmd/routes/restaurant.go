package routes

import (
	"venues/cmd/controllers"

	"github.com/labstack/echo"
)

func BuildRestaurantGroup(group *echo.Group) {
	controller := controllers.NewRestaurantController()
	group.GET("", controller.List)
	group.POST("", controller.Create)
	group.POST("/:restaurant_id", controller.Update)
	group.DELETE("/:restaurant_id", controller.Remove)
	group.POST("/:restaurant_id/dish", controller.AddDish)
}
