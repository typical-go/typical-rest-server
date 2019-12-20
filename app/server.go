package app

import (
	"fmt"

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
	controller.MusicCntrl
}

func (s server) Middleware() {
	s.Use(middleware.Recover())
}

func (s server) Route() {
	s.AppCntrl.Route(s.Echo)
	s.BookCntrl.Route(s.Echo)
	s.MusicCntrl.Route(s.Echo)
}

func (s server) Start() (err error) {
	if err = s.Echo.Start(s.Config.Address); err != nil {
		return fmt.Errorf("App: %w", err)
	}
	return
}
