package typserver

import (
	"github.com/labstack/echo"
)

type Controller interface {
	SetRoute(*echo.Echo)
}
