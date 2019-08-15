package typserver

import (
	"context"
	"time"

	logrusmiddleware "github.com/bakatz/echo-logrusmiddleware"
	"github.com/labstack/echo"
	log "github.com/sirupsen/logrus"
)

var server *echo.Echo

// Create new server
func Create(cfg *Config) *echo.Echo {
	if server == nil {
		log.Infof("Create http server; DEBUG=%t", cfg.Debug)

		server = echo.New()
		server.HideBanner = true
		server.Debug = cfg.Debug

		logrusMwConfig := logrusmiddleware.Config{}

		// set log
		if cfg.Debug {
			log.SetLevel(log.InfoLevel)
			logrusMwConfig.IncludeRequestBodies = true
			logrusMwConfig.IncludeResponseBodies = true
		} else {
			log.SetLevel(log.WarnLevel)
			log.SetFormatter(&log.JSONFormatter{})
		}

		server.Use(logrusmiddleware.HookWithConfig(logrusMwConfig))
		server.Logger = logrusmiddleware.Logger{Logger: log.StandardLogger()}
	}

	return server
}

// Shutdown the server
func Shutdown(server *echo.Echo) error {
	log.Info("Server is shutting down")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	return server.Shutdown(ctx)
}
