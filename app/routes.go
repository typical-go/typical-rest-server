package app

import (
	"net/http"

	"github.com/labstack/echo"
)

func (s *server) initRoutes() {
	s.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})
}
