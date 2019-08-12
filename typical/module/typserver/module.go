package typserver

import (
	"context"
	"time"

	logrusmiddleware "github.com/bakatz/echo-logrusmiddleware"
	"github.com/labstack/echo"
	log "github.com/sirupsen/logrus"
	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/typictx"
)

// Module for postgres
func Module(logHooks ...log.Hook) *typictx.Module {
	return &typictx.Module{
		Name:         "Echo Server with Logrus",
		ConfigPrefix: "SERVER",
		ConfigSpec:   &Config{},
		OpenFunc: func(cfg *Config) (server *echo.Echo, err error) {
			server = echo.New()
			server.HideBanner = true
			server.Debug = cfg.Debug

			logrusMwConfig := logrusmiddleware.Config{}

			// set log
			if cfg.Debug {
				log.SetLevel(log.InfoLevel)
				logrusMwConfig.IncludeRequestBodies = true
				logrusMwConfig.IncludeResponseBodies = true
				log.Info("Configure logger for debug mode")
			} else {
				log.SetLevel(log.WarnLevel)
				log.SetFormatter(&log.JSONFormatter{})
			}

			for _, hook := range logHooks {
				log.AddHook(hook)
			}

			server.Use(logrusmiddleware.HookWithConfig(logrusMwConfig))
			server.Logger = logrusmiddleware.Logger{Logger: log.StandardLogger()}
			return
		},
		CloseFunc: func(server *echo.Echo) error {
			log.Info("Server is shutting down")
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()
			return server.Shutdown(ctx)
		},
	}
}
