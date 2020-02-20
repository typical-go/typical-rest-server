package restserver

import (
	"context"
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"

	logrusmiddleware "github.com/bakatz/echo-logrusmiddleware"
	"github.com/labstack/echo"
	"github.com/typical-go/typical-go/pkg/typapp"
	"github.com/typical-go/typical-go/pkg/typcfg"
	"github.com/typical-go/typical-rest-server/restserver/config"
	"github.com/urfave/cli/v2"
)

// Module of application
type Module struct {
	prefix string
	Debug  bool
}

// New application [nowire]
func New() *Module {
	return &Module{
		prefix: "APP",
	}
}

// WithPrefix return module with new prefix
func (m *Module) WithPrefix(prefix string) *Module {
	m.prefix = prefix
	return m
}

// EntryPoint of application
func (m *Module) EntryPoint() interface{} {
	return func(s server) error {
		s.Middleware()
		s.Route()
		return s.Start()
	}
}

// AppCommands return comamnd
func (m *Module) AppCommands(c *typapp.Context) []*cli.Command {
	return []*cli.Command{
		{
			Name:   "route",
			Usage:  "Print available API Routes",
			Action: c.ActionFunc(taskRouteList),
		},
	}
}

// WithDebug to set debug
func (m *Module) WithDebug(debug bool) *Module {
	m.Debug = debug
	return m
}

// Configure server
func (m *Module) Configure(loader typcfg.Loader) *typcfg.Detail {
	return &typcfg.Detail{
		Prefix: m.prefix,
		Spec: &config.Config{
			Debug: m.Debug,
		},
		Constructor: func() (cfg config.Config, err error) {
			err = loader.Load(m.prefix, &cfg)
			return
		},
	}
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
// TODO: create helper
func (m *Module) Create(cfg config.Config) *echo.Echo {
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
// TODO: create helper
func (m *Module) Shutdown(server *echo.Echo) (err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err = server.Shutdown(ctx); err != nil {
		return fmt.Errorf("Server: Shutdown: %w", err)
	}
	return
}
