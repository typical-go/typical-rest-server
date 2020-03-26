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

// Module of new-relic
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

// Configure the module
func (m *Module) Configure() *typcfg.Configuration {
	return typcfg.NewConfiguration(m.configName, &Config{})
}

// Provide dependencies
func (*Module) Provide() []*typapp.Constructor {
	return []*typapp.Constructor{
		typapp.NewConstructor(func(cfg *Config) (newrelic.Application, error) {
			if cfg.AppName == "" || cfg.Key == "" {
				return nil, nil
			}
			nrCfg := newrelic.NewConfig(cfg.AppName, cfg.Key)
			return newrelic.NewApplication(nrCfg)
		}),
	}
}
