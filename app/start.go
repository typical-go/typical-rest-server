package app

import (
	"github.com/labstack/echo"
	log "github.com/sirupsen/logrus"
	"github.com/typical-go/typical-rest-server/app/config"
)

// Start the service
func Start(server *echo.Echo, cfg *config.Config) error {
	log.Info("Start the application")
	return server.Start(cfg.Address)
}
