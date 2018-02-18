package assembly

import (
	"net/http"

	"venues/pkg/healthcheckers"

	"github.com/labstack/echo"
)

type HealthCheck struct {
	checkers []healthcheckers.Checker
}

func (h *HealthCheck) Check(c echo.Context) error {
	for _, checker := range h.checkers {
		if err := checker.Check(); err != nil {
			return c.String(http.StatusServiceUnavailable, checker.Message())
		}
	}

	return c.String(http.StatusOK, "Ok")
}
