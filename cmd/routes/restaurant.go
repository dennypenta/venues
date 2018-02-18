package routes

import (
	"venues/cmd/controllers"

	"github.com/labstack/echo"
)

func BuildRestaurantGroup(group *echo.Group) {
	controller := controllers.NewRestaurantController()
	group.GET("/restaurants", controller.List)
	group.POST("/restaurants", controller.Create)
	group.POST("/restaurants/:restaurant_id", controller.Update)
	group.DELETE("/restaurants/:restaurant_id", controller.Remove)
}
