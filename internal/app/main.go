package app

import (
	"github.com/labstack/echo/v4"
)

// Main function to run server
func Main(a app) (err error) {
	e := echo.New()
	defer shutdown(e)

	e.HideBanner = true
	e.Debug = a.Debug

	setMiddleware(e)

	if err := setRoute(e, &a); err != nil {
		return err
	}

	return e.Start(a.Address)
}
