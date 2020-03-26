package server

import (
	"github.com/typical-go/typical-go/pkg/typapp"
	"github.com/typical-go/typical-go/pkg/typcfg"
	"github.com/typical-go/typical-rest-server/pkg/typserver"
	"github.com/typical-go/typical-rest-server/server/config"
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
func (m *App) Configure() *typcfg.Configuration {
	return typcfg.NewConfiguration(m.configName, &config.Config{
		Debug: m.debug,
	})
}

// EntryPoint of application
func (m *App) EntryPoint() *typapp.MainInvocation {
	return typapp.NewMainInvocation(startServer)
}

// Provide dependencies
func (m *App) Provide() []*typapp.Constructor {
	return []*typapp.Constructor{
		typapp.NewConstructor(typserver.New),
	}
}

// Destroy dependencies
func (m *App) Destroy() []*typapp.Destruction {
	return []*typapp.Destruction{
		typapp.NewDestruction(typserver.Shutdown),
	}
}
