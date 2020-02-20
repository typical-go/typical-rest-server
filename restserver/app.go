package restserver

import (
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/typical-go/typical-rest-server/restserver/config"
	"github.com/typical-go/typical-rest-server/restserver/controller"

	"go.uber.org/dig"
)

type app struct {
	dig.In
	*echo.Echo
	config.Config
	controller.BookCntrl
	controller.AppCntrl
}

func (a app) Middleware() {
	a.Use(middleware.Recover())
}

func (a app) Route() {
	a.AppCntrl.Route(a.Echo)
	a.BookCntrl.Route(a.Echo)
}

func (a app) Start() (err error) {
	a.Middleware()
	a.Route()
	return a.Echo.Start(a.Config.Address)
}
