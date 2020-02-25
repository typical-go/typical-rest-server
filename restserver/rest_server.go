package restserver

import (
	"github.com/labstack/echo"
	"github.com/typical-go/typical-go/pkg/typcfg"
	"github.com/typical-go/typical-go/pkg/typdep"
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

// Configure server
func (m *RestServer) Configure(loader typcfg.Loader) *typcfg.Detail {
	return &typcfg.Detail{
		Prefix: m.prefix,
		Spec: &config.Config{
			Debug: m.debug,
		},
		Constructor: typdep.NewConstructor(
			func() (cfg config.Config, err error) {
				err = loader.Load(m.prefix, &cfg)
				return
			}),
	}
}

// EntryPoint of application
func (m *RestServer) EntryPoint() *typdep.Invocation {
	return typdep.NewInvocation(func(a app) error {
		return a.Start()
	})
}

// Provide dependencies
func (m *RestServer) Provide() []*typdep.Constructor {
	return []*typdep.Constructor{
		typdep.NewConstructor(func(cfg config.Config) *echo.Echo {
			return serverkit.Create(cfg.Debug)
		}),
	}
}

// Destroy dependencies
func (m *RestServer) Destroy() []*typdep.Invocation {
	return []*typdep.Invocation{
		typdep.NewInvocation(serverkit.Shutdown),
	}
}
