package restserver

import (
	"github.com/labstack/echo/middleware"
	"github.com/typical-go/typical-rest-server/pkg/typserver"
	"github.com/typical-go/typical-rest-server/restserver/config"
	"github.com/typical-go/typical-rest-server/restserver/controller"

	"go.uber.org/dig"
)

type app struct {
	dig.In
	*typserver.Server
	config.Config
	controller.BookCntrl
	controller.AppCntrl
}

func startServer(a app) (err error) {
	a.SetDebug(a.Debug)
	a.Use(middleware.Recover())
	a.AppCntrl.Route(a.Echo)
	a.BookCntrl.Route(a.Echo)
	return a.Start(a.Address)
}
