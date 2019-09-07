package app

import (
	"github.com/labstack/echo"
	log "github.com/sirupsen/logrus"
	"github.com/typical-go/typical-rest-server/app/controller"
	"go.uber.org/dig"
)

// RouteParams for route parameters
type RouteParams struct {
	dig.In
	*echo.Echo
	controller.BookController
	controller.ApplicationController
}

// Routes of API
func Routes(p RouteParams) {
	log.Info("Initiate API Routes")
	p.BookController.Route(p.Echo)
	p.BookController.Route(p.Echo)
}
