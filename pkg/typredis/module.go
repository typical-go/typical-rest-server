package typredis

import (
	"github.com/typical-go/typical-go/pkg/common"
	"github.com/typical-go/typical-go/pkg/typcore"
	"github.com/urfave/cli/v2"
)

// Module of Redis
type Module struct {
	Host     string
	Port     string
	Password string
}

// New Redis Module
func New() *Module {
	return &Module{
		Host:     "localhost",
		Port:     "6379",
		Password: "redispass",
	}
}

// WithHost to set host
func (m *Module) WithHost(host string) *Module {
	m.Host = host
	return m
}

// WithPort to set port
func (m *Module) WithPort(port string) *Module {
	m.Port = port
	return m
}

// WithPassword to set password
func (m *Module) WithPassword(password string) *Module {
	m.Password = password
	return m
}

// BuildCommands of module
func (m *Module) BuildCommands(c *typcore.BuildContext) []*cli.Command {
	return []*cli.Command{
		{
			Name:  "redis",
			Usage: "Redis Tool",
			Before: func(ctx *cli.Context) error {
				return common.LoadEnvFile()
			},
			Subcommands: []*cli.Command{
				m.consoleCmd(c),
			},
		},
	}
}

// Configure Redis
func (m *Module) Configure(loader typcore.ConfigLoader) (prefix string, spec, loadFn interface{}) {
	prefix = "REDIS"
	spec = &Config{
		Host:     m.Host,
		Port:     m.Port,
		Password: m.Password,
	}
	loadFn = func() (cfg Config, err error) {
		err = loader.Load(prefix, &cfg)
		return
	}
	return
}

// Provide dependencies
func (m *Module) Provide() []interface{} {
	return []interface{}{
		m.connect,
	}
}

// Prepare the module
func (m *Module) Prepare() []interface{} {
	return []interface{}{
		m.ping,
	}
}

// Destroy dependencies
func (m *Module) Destroy() []interface{} {
	return []interface{}{
		m.disconnect,
	}
}
