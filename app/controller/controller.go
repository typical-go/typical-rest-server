package controller

import (
	"net/http"

	"github.com/labstack/echo"
)

func underContruction(e echo.Context) error {
	return e.JSON(http.StatusServiceUnavailable, map[string]string{"message": "Under Construction"})
}
