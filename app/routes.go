package app

import (
	"github.com/labstack/echo"
	log "github.com/sirupsen/logrus"
	"github.com/typical-go/typical-rest-server/app/controller"
	"go.uber.org/dig"
)

// Controllers is collection fo controller
type Controllers struct {
	dig.In

	Book *controller.BookController
	App  *controller.ApplicationController
}

// Routes of API
func Routes(server *echo.Echo, c Controllers) {
	log.Info("Initiate API Routes")

	c.Book.Route(server)
	c.App.Route(server)
}
