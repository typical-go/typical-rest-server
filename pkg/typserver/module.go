package typserver

import (
	"context"
	"fmt"
	"time"

	logrusmiddleware "github.com/bakatz/echo-logrusmiddleware"
	"github.com/labstack/echo"
	log "github.com/sirupsen/logrus"
	"github.com/typical-go/typical-go/pkg/typcore"
)

// Module of Server
func Module() interface{} {
	return &module{}
}

type module struct{}

// Configure server
func (s *module) Configure() (prefix string, spec, loadFn interface{}) {
	prefix = "SERVER"
	spec = &Config{}
	loadFn = func(loader typcore.ConfigLoader) (cfg Config, err error) {
		err = loader.Load(prefix, &cfg)
		return
	}
	return
}

// Provide dependencies
func (s *module) Provide() []interface{} {
	return []interface{}{
		s.Create,
	}
}

// Destroy dependencies
func (s *module) Destroy() []interface{} {
	return []interface{}{
		s.Shutdown,
	}
}

// Create new server
func (s *module) Create(cfg Config) *echo.Echo {
	server := echo.New()
	server.HideBanner = true
	server.Debug = cfg.Debug
	logrusMwConfig := logrusmiddleware.Config{}
	if cfg.Debug {
		log.SetLevel(log.DebugLevel)
		logrusMwConfig.IncludeRequestBodies = true
		logrusMwConfig.IncludeResponseBodies = true
	} else {
		log.SetLevel(log.WarnLevel)
		log.SetFormatter(&log.JSONFormatter{})
	}
	server.Use(logrusmiddleware.HookWithConfig(logrusMwConfig))
	server.Logger = logrusmiddleware.Logger{Logger: log.StandardLogger()}
	return server
}

// Shutdown the server
func (s *module) Shutdown(server *echo.Echo) error {
	fmt.Println("Server is shutting down")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	return server.Shutdown(ctx)
}
