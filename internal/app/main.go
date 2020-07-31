package app

import (
	"log"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"github.com/typical-go/typical-rest-server/pkg/echokit"
)

// Main function to run server
func Main(a app) (err error) {
	e := echo.New()
	defer shutdown(e)

	e.HideBanner = true
	e.Debug = a.Debug

	logger := logrus.StandardLogger()
	setLogLevel(logger, a.Debug)
	e.Logger = echokit.WrapLogrus(logger) // NOTE: setup echo logger
	log.SetOutput(logger.Writer())        // NOTE: setup golang std log

	setMiddleware(e)

	if err := setRoute(e, &a); err != nil {
		return err
	}
	return e.Start(a.Address)
}
