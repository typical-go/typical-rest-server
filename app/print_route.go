package app

import (
	"fmt"

	"github.com/labstack/echo"
	"github.com/typical-go/typical-rest-server/app/config"
)

// Route to return routes
func Route(e *echo.Echo, cfg *config.Config) (err error) {
	for _, route := range e.Routes() {
		fmt.Println(route.Path)
	}
	return
}
