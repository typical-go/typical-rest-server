package typredis

import (
	"github.com/typical-go/typical-go/pkg/typbuild"
	"github.com/typical-go/typical-go/pkg/typcfg"
	"github.com/typical-go/typical-go/pkg/typcore"
	"github.com/urfave/cli/v2"
)

// Module of Redis
type Module struct {
	Host        string
	Port        string
	Password    string
	DockerName  string
	DockerImage string
}

// New Redis Module
func New() *Module {
	return &Module{
		Host:        "localhost",
		Port:        "6379",
		Password:    "redispass",
		DockerImage: "redis:4.0.5-alpine",
		DockerName:  "redis",
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

// WithDockerImage to set docker image
func (m *Module) WithDockerImage(dockerImage string) *Module {
	m.DockerImage = dockerImage
	return m
}

// WithDockerName to set docker name
func (m *Module) WithDockerName(dockerName string) *Module {
	m.DockerName = dockerName
	return m
}

// BuildCommands of module
func (m *Module) BuildCommands(c *typbuild.Context) []*cli.Command {
	return []*cli.Command{
		{
			Name:  "redis",
			Usage: "Redis Tool",
			Before: func(ctx *cli.Context) error {
				return typcfg.LoadEnvFile()
			},
			Subcommands: []*cli.Command{
				{
					Name:    "console",
					Aliases: []string{"c"},
					Usage:   "Redis Interactive",
					Action:  c.ActionFunc(m.console),
				},
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
