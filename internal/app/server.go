package app

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/typical-go/typical-rest-server/internal/app/controller"
	"github.com/typical-go/typical-rest-server/internal/app/infra"
	"github.com/typical-go/typical-rest-server/pkg/echokit"
)

// SetServer to echo server
func SetServer(
	e *echo.Echo,
	bookCntrl controller.BookCntrl,
	songCntrl controller.SongCntrl,
) {

	// set middleware
	e.Use(infra.LogMiddleware)
	e.Use(middleware.Recover())

	// set route
	echokit.SetRoute(e,
		&bookCntrl,
		&songCntrl,
	)
}
