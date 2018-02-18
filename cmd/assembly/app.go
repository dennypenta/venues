package assembly

import (
	"fmt"
	"net/http"

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
	app.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!\n")
	})
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
