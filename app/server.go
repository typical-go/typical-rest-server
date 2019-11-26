package app

import (
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/typical-go/typical-rest-server/app/config"
	"github.com/typical-go/typical-rest-server/app/controller"
	"go.uber.org/dig"
)

type server struct {
	dig.In
	*echo.Echo
	config.Config
	controller.BookCntrl
	controller.AppCntrl
}

// Start the service
func startServer(p server) error {
	// Middlewares
	p.Use(middleware.Recover())

	// Routes
	p.AppCntrl.Route(p.Echo)
	p.BookCntrl.Route(p.Echo)

	return p.Echo.Start(p.Config.Address)
}
