package typserver

import (
	"context"
	"fmt"
	"time"

	"github.com/typical-go/typical-go/pkg/typmod"

	logrusmiddleware "github.com/bakatz/echo-logrusmiddleware"
	"github.com/labstack/echo"
	log "github.com/sirupsen/logrus"
)

// Config is server configuration
type Config struct {
	Debug bool `default:"false"`
}

// Module of server
func Module() interface{} {
	return &serverModule{
		Name: "Server",
		Configuration: typmod.Configuration{
			Prefix: "SERVER",
			Spec:   &Config{},
		},
	}
}

type serverModule struct {
	typmod.Configuration
	Name string
}

func (s serverModule) Provide() []interface{} {
	return []interface{}{
		s.loadConfig,
		s.Create,
	}
}

func (s serverModule) Destroy() []interface{} {
	return []interface{}{
		s.Shutdown,
	}
}

func (s serverModule) loadConfig() (cfg *Config, err error) {
	err = s.Configuration.Load()
	cfg = s.Configuration.Spec.(*Config)
	return
}

// Create new server
func (s serverModule) Create(cfg *Config) *echo.Echo {
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
func (s serverModule) Shutdown(server *echo.Echo) error {
	fmt.Println("Server is shutting down")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	return server.Shutdown(ctx)
}
