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
	Server    *echo.Echo
	BookCntrl controller.BookController
	AppCntrl  controller.ApplicationController
}

// Routes of API
func Routes(p RouteParams) {
	log.Info("Initiate API Routes")
	p.BookCntrl.Route(p.Server)
	p.AppCntrl.Route(p.Server)
}
