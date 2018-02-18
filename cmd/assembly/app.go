package assembly

import (
	"fmt"

	"venues/cmd/routes"
	"venues/cmd/storages"
	"venues/pkg/healthcheckers"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

type App struct {
	*echo.Echo
}

func (app *App) setMiddleware() {
	app.Use(middleware.Logger())
	app.Use(middleware.Recover())
}

func (app *App) setRoutes() {
	mongoHealthChecker := &healthcheckers.CheckService{
		ServiceName: "Mongo",
		Action:      storages.GetStorage().Session.Ping,
	}
	healthCkecker := HealthCheck{[]healthcheckers.Checker{mongoHealthChecker}}
	app.GET("/", healthCkecker.Check)

	restaurantGroup := app.Group("/restaurants")
	routes.BuildRestaurantGroup(restaurantGroup)
}

func (app *App) init() {
	app.setMiddleware()

	app.setRoutes()
}

func (app *App) Run(port string) {
	startPort := fmt.Sprintf(":%s", port)
	app.Logger.Fatal(app.Start(startPort))
}

func NewApp() *App {
	app := &App{echo.New()}
	app.init()

	return app
}
