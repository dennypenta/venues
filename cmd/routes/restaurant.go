package routes

import (
	"venues/cmd/controllers"

	"github.com/labstack/echo"
)

func BuildRestaurantGroup(group *echo.Group) {
	controller := controllers.NewRestaurantController()
	group.GET("/restaurants", controller.List)
}
