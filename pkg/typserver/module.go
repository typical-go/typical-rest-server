package typserver

import (
	"context"
	"fmt"
	"time"

	logrusmiddleware "github.com/bakatz/echo-logrusmiddleware"
	"github.com/labstack/echo"
	log "github.com/sirupsen/logrus"
	"github.com/typical-go/typical-go/pkg/typcfg"
)

// Module of server
type Module struct {
	Debug bool
}

// New server module
func New() *Module {
	return &Module{}
}

// WithDebug to set debug
func (m *Module) WithDebug(debug bool) *Module {
	m.Debug = debug
	return m
}

// Configure server
func (m *Module) Configure(loader typcfg.Loader) (prefix string, spec, constructor interface{}) {
	prefix = "SERVER"
	spec = &Config{
		Debug: m.Debug,
	}
	constructor = func() (cfg Config, err error) {
		err = loader.Load(prefix, &cfg)
		return
	}
	return
}

// Provide dependencies
func (m *Module) Provide() []interface{} {
	return []interface{}{
		m.Create,
	}
}

// Destroy dependencies
func (m *Module) Destroy() []interface{} {
	return []interface{}{
		m.Shutdown,
	}
}

// Create new server
func (m *Module) Create(cfg Config) *echo.Echo {
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
func (m *Module) Shutdown(server *echo.Echo) (err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err = server.Shutdown(ctx); err != nil {
		return fmt.Errorf("Server: Shutdown: %w", err)
	}
	return
}
