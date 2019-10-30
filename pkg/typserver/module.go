package typserver

import (
	"context"
	"time"

	"github.com/kelseyhightower/envconfig"
	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/typiobj"
	"go.uber.org/dig"

	logrusmiddleware "github.com/bakatz/echo-logrusmiddleware"
	"github.com/labstack/echo"
	log "github.com/sirupsen/logrus"
)

// Module of server
func Module() interface{} {
	return &serverModule{
		Name: "Server",
		Configuration: typiobj.Configuration{
			Prefix: "SERVER",
			Spec:   &Config{},
		},
	}
}

type serverModule struct {
	typiobj.Configuration
	Name string
}

func (s serverModule) Provide() []interface{} {
	return []interface{}{
		s.loadConfig,
		s.Create,
	}
}

func (s serverModule) Destruct(c *dig.Container) (err error) {
	return c.Invoke(s.Shutdown)
}

func (s serverModule) loadConfig() (cfg *Config, err error) {
	cfg = new(Config)
	err = envconfig.Process(s.Configure().Prefix, cfg)
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
	log.Info("Server is shutting down")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	return server.Shutdown(ctx)
}
