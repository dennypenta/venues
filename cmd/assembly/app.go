package assembly

import (
	"fmt"

	"venues/cmd/routes"
	"venues/cmd/storages"
	"venues/pkg/healthcheckers"
	"venues/pkg/validator"

	"context"
	"os"
	"os/signal"
	"time"

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
	go app.Logger.Fatal(app.Start(startPort))

	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	// do not forget about bubble ordering of defers
	// close sessions before
	defer cancel()
	defer storages.GetStorage().Session.Close()
	if err := app.Shutdown(ctx); err != nil {
		app.Logger.Fatal(err)
	}

}

func NewApp() *App {
	app := &App{echo.New()}

	// setup validator that will be used by echo.Context.Bind
	app.Validator = validator.NewValidator()

	app.init()

	return app
}
