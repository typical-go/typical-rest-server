package infra

import (
	"github.com/labstack/echo/v4"
	"github.com/typical-go/typical-rest-server/internal/app/infra/log"
	"github.com/typical-go/typical-rest-server/pkg/logruskit"
)

// NewServer return new instance of server
// @ctor
func NewServer(cfg *AppCfg) *echo.Echo {
	e := echo.New()
	logger := log.SetDebug(cfg.Debug)

	e.HideBanner = true
	e.Debug = cfg.Debug
	e.Logger = logruskit.EchoLogger(logger)
	return e
}
