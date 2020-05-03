package typnewrelic

import (
	newrelic "github.com/newrelic/go-agent"
	"github.com/typical-go/typical-go/pkg/typapp"
	"github.com/typical-go/typical-go/pkg/typcfg"
)

const (
	// DefaultConfigName is default value for ConfigName
	DefaultConfigName = "NEWRELIC"
)

var (
	_ typapp.Provider   = (*Module)(nil)
	_ typcfg.Configurer = (*Module)(nil)
)

// Module of new-relic
// TODO: use typapp.Module instead
type Module struct {
	configName string
}

// New instance of new-relic module
func New() *Module {
	return &Module{
		configName: DefaultConfigName,
	}
}

// WithConfigName return new-relic module with new config-name
func (m *Module) WithConfigName(configName string) *Module {
	m.configName = configName
	return m
}

// Configurations the module
func (m *Module) Configurations() []*typcfg.Configuration {
	return []*typcfg.Configuration{
		&typcfg.Configuration{
			Name: m.configName,
			Spec: &Config{},
		},
	}
}

// Constructors of newreslice
func (*Module) Constructors() []*typapp.Constructor {
	return []*typapp.Constructor{
		typapp.NewConstructor("", createApp),
	}
}

func createApp(cfg *Config) (newrelic.Application, error) {
	if cfg.AppName == "" || cfg.Key == "" {
		return nil, nil
	}
	nrCfg := newrelic.NewConfig(cfg.AppName, cfg.Key)
	return newrelic.NewApplication(nrCfg)
}
