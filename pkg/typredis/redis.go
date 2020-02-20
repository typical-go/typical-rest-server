package typredis

import (
	"github.com/typical-go/typical-go/pkg/typbuild"
	"github.com/typical-go/typical-go/pkg/typcfg"
	"github.com/urfave/cli/v2"
)

// Redis of Redis
type Redis struct {
	Host        string
	Port        string
	Password    string
	DockerName  string
	DockerImage string
	prefix      string
}

// New Redis Redis
func New() *Redis {
	return &Redis{
		Host:        "localhost",
		Port:        "6379",
		Password:    "redispass",
		DockerImage: "redis:4.0.5-alpine",
		DockerName:  "redis",
		prefix:      "REDIS",
	}
}

// WithHost to return module with new host
func (m *Redis) WithHost(host string) *Redis {
	m.Host = host
	return m
}

// WithPort to return module with new port
func (m *Redis) WithPort(port string) *Redis {
	m.Port = port
	return m
}

// WithPassword to return module with new password
func (m *Redis) WithPassword(password string) *Redis {
	m.Password = password
	return m
}

// WithDockerImage to return module with new docker image
func (m *Redis) WithDockerImage(dockerImage string) *Redis {
	m.DockerImage = dockerImage
	return m
}

// WithDockerName to return module with new docker name
func (m *Redis) WithDockerName(dockerName string) *Redis {
	m.DockerName = dockerName
	return m
}

// WithPrefix to return module with new prefix
func (m *Redis) WithPrefix(prefix string) *Redis {
	m.prefix = prefix
	return m
}

// BuildCommands of module
func (m *Redis) BuildCommands(c *typbuild.Context) []*cli.Command {
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
func (m *Redis) Configure(loader typcfg.Loader) *typcfg.Detail {
	return &typcfg.Detail{
		Prefix: m.prefix,
		Spec: &Config{
			Host:     m.Host,
			Port:     m.Port,
			Password: m.Password,
		},
		Constructor: func() (cfg Config, err error) {
			err = loader.Load(m.prefix, &cfg)
			return
		},
	}
}

// Provide dependencies
func (m *Redis) Provide() []interface{} {
	return []interface{}{
		m.connect,
	}
}

// Prepare the module
func (m *Redis) Prepare() []interface{} {
	return []interface{}{
		m.ping,
	}
}

// Destroy dependencies
func (m *Redis) Destroy() []interface{} {
	return []interface{}{
		m.disconnect,
	}
}
