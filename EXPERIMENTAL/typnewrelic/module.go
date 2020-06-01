package typnewrelic

import (
	newrelic "github.com/newrelic/go-agent"
	"github.com/typical-go/typical-go/pkg/typgo"
)

var (
	// DefaultConfigName is default value for ConfigName
	DefaultConfigName = "NEWRELIC"
)

var (
	_ typgo.Provider   = (*Module)(nil)
	_ typgo.Configurer = (*Module)(nil)
)

// Module of new-relic
// TODO: use typgo.Module instead
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
func (m *Module) Configurations() []*typgo.Configuration {
	return []*typgo.Configuration{
		&typgo.Configuration{
			Name: m.configName,
			Spec: &Config{},
		},
	}
}

// Constructors of newreslice
func (*Module) Constructors() []*typgo.Constructor {
	return []*typgo.Constructor{
		&typgo.Constructor{Fn: createApp},
	}
}

func createApp(cfg *Config) (newrelic.Application, error) {
	if cfg.AppName == "" || cfg.Key == "" {
		return nil, nil
	}
	nrCfg := newrelic.NewConfig(cfg.AppName, cfg.Key)
	return newrelic.NewApplication(nrCfg)
}
