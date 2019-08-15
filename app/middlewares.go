package app

import (
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	log "github.com/sirupsen/logrus"
)

// Middlewares for the service
func Middlewares(server *echo.Echo) {
	log.Info("Initiate Middlewares")

	server.Use(middleware.Recover())
}
