package app

import (
	"log"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"github.com/typical-go/typical-rest-server/pkg/typrest"
)

// Main function to run server
func Main(a app) (err error) {
	e := echo.New()
	defer shutdown(e)

	e.HideBanner = true
	e.Debug = a.Debug

	logger := logrus.StandardLogger() // NOTE: always use standard logrus logger
	setLogLevel(logger, a.Debug)
	e.Logger = typrest.WrapLogrus(logger) // NOTE: setup echo logger
	log.SetOutput(logger.Writer())        // NOTE: std golang log use same output writer with logrus

	setMiddleware(e)

	if err := setRoute(e, &a); err != nil {
		return err
	}
	return e.Start(a.Address)
}
