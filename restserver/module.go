package restserver

import (
	"github.com/labstack/echo"
	"github.com/typical-go/typical-go/pkg/typcfg"
	"github.com/typical-go/typical-rest-server/pkg/serverkit"
	"github.com/typical-go/typical-rest-server/restserver/config"
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
		func(cfg config.Config) *echo.Echo {
			return serverkit.Create(cfg.Debug)
		},
	}
}

// Destroy dependencies
func (m *Module) Destroy() []interface{} {
	return []interface{}{
		serverkit.Shutdown,
	}
}
