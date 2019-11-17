package app

import (
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/typical-go/typical-rest-server/app/config"
	"github.com/typical-go/typical-rest-server/app/controller"
	"go.uber.org/dig"
)

type params struct {
	dig.In
	*echo.Echo
	config.Config
	controller.BookCntrl
	controller.AppCntrl
}

// Routes of API
func Routes(p params) {
	p.AppCntrl.Route(p.Echo)
	p.BookCntrl.Route(p.Echo)
}

// Middlewares for the service
func Middlewares(p params) {
	p.Echo.Use(middleware.Recover())
}

// Start the service
func Start(p params) error {
	return p.Echo.Start(p.Config.Address)
}
