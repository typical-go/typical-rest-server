package infra

import (
	"github.com/labstack/echo/v4"
	"github.com/typical-go/typical-rest-server/pkg/logruskit"
)

// NewEcho return new instance of server
// @ctor
func NewEcho(cfg *AppCfg) *echo.Echo {
	e := echo.New()
	logger := SetLogger(cfg.Debug)

	e.HideBanner = true
	e.Debug = cfg.Debug
	e.Logger = logruskit.EchoLogger(logger)
	return e
}
