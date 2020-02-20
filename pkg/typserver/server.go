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

// Server of server
type Server struct {
	Debug  bool
	prefix string
}

// New server module
func New() *Server {
	return &Server{
		prefix: "SERVER",
	}
}

// WithDebug to set debug
func (m *Server) WithDebug(debug bool) *Server {
	m.Debug = debug
	return m
}

// Configure server
func (m *Server) Configure(loader typcfg.Loader) *typcfg.Detail {
	return &typcfg.Detail{
		Prefix: m.prefix,
		Spec: &Config{
			Debug: m.Debug,
		},
		Constructor: func() (cfg Config, err error) {
			err = loader.Load(m.prefix, &cfg)
			return
		},
	}
}

// Provide dependencies
func (m *Server) Provide() []interface{} {
	return []interface{}{
		m.Create,
	}
}

// Destroy dependencies
func (m *Server) Destroy() []interface{} {
	return []interface{}{
		m.Shutdown,
	}
}

// Create new server
func (m *Server) Create(cfg Config) *echo.Echo {
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
func (m *Server) Shutdown(server *echo.Echo) (err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err = server.Shutdown(ctx); err != nil {
		return fmt.Errorf("Server: Shutdown: %w", err)
	}
	return
}
