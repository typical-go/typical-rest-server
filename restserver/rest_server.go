package restserver

import (
	"github.com/labstack/echo"
	"github.com/typical-go/typical-go/pkg/typcfg"
	"github.com/typical-go/typical-rest-server/pkg/serverkit"
	"github.com/typical-go/typical-rest-server/restserver/config"
)

// RestServer of application
type RestServer struct {
	prefix string
	debug  bool
}

// New application [nowire]
func New() *RestServer {
	return &RestServer{
		prefix: "APP",
	}
}

// WithPrefix return module with new prefix
func (m *RestServer) WithPrefix(prefix string) *RestServer {
	m.prefix = prefix
	return m
}

// WithDebug to set debug
func (m *RestServer) WithDebug(debug bool) *RestServer {
	m.debug = debug
	return m
}

// EntryPoint of application
func (m *RestServer) EntryPoint() interface{} {
	return func(a app) error {
		return a.Start()
	}
}

// Configure server
func (m *RestServer) Configure(loader typcfg.Loader) *typcfg.Detail {
	return &typcfg.Detail{
		Prefix: m.prefix,
		Spec: &config.Config{
			Debug: m.debug,
		},
		Constructor: func() (cfg config.Config, err error) {
			err = loader.Load(m.prefix, &cfg)
			return
		},
	}
}

// Provide dependencies
func (m *RestServer) Provide() []interface{} {
	return []interface{}{
		func(cfg config.Config) *echo.Echo {
			return serverkit.Create(cfg.Debug)
		},
	}
}

// Destroy dependencies
func (m *RestServer) Destroy() []interface{} {
	return []interface{}{
		serverkit.Shutdown,
	}
}
