package restserver

import (
	"github.com/typical-go/typical-go/pkg/typcfg"
	"github.com/typical-go/typical-go/pkg/typdep"
	"github.com/typical-go/typical-rest-server/pkg/typserver"
	"github.com/typical-go/typical-rest-server/restserver/config"
)

// App of rest server
type App struct {
	configName string
	debug      bool
}

// New application [nowire]
func New() *App {
	return &App{
		configName: "APP",
	}
}

// WithConfigName return module with new prefix
func (m *App) WithConfigName(configName string) *App {
	m.configName = configName
	return m
}

// WithDebug to set debug
func (m *App) WithDebug(debug bool) *App {
	m.debug = debug
	return m
}

// Configure server
func (m *App) Configure(loader typcfg.Loader) *typcfg.Configuration {
	return &typcfg.Configuration{
		Name: m.configName,
		Spec: &config.Config{
			Debug: m.debug,
		},
		Constructor: typdep.NewConstructor(
			func() (cfg config.Config, err error) {
				err = loader.Load(m.configName, &cfg)
				return
			}),
	}
}

// EntryPoint of application
func (m *App) EntryPoint() *typdep.Invocation {
	return typdep.NewInvocation(startServer)
}

// Provide dependencies
func (m *App) Provide() []*typdep.Constructor {
	return []*typdep.Constructor{
		typdep.NewConstructor(typserver.New),
	}
}

// Destroy dependencies
func (m *App) Destroy() []*typdep.Invocation {
	return []*typdep.Invocation{
		typdep.NewInvocation(typserver.Shutdown),
	}
}
