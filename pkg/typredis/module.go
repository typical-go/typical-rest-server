package typredis

import (
	"github.com/typical-go/typical-go/pkg/typcore"
	"github.com/typical-go/typical-go/pkg/common"
	"github.com/urfave/cli/v2"
)

// Module of redis
func Module() interface{} {
	return &module{}
}

type module struct{}

// BuildCommands of module
func (r *module) BuildCommands(c *typcore.Context) []*cli.Command {
	return []*cli.Command{
		{
			Name:  "redis",
			Usage: "Redis Tool",
			Before: func(ctx *cli.Context) error {
				return common.LoadEnvFile()
			},
			Subcommands: []*cli.Command{
				r.consoleCmd(c),
			},
		},
	}
}

// Configure Redis
func (r *module) Configure() (prefix string, spec, loadFn interface{}) {
	prefix = "REDIS"
	spec = &Config{}
	loadFn = func(loader typcore.ConfigLoader) (cfg Config, err error) {
		err = loader.Load(prefix, &cfg)
		return
	}
	return
}

// Provide dependencies
func (r *module) Provide() []interface{} {
	return []interface{}{
		r.connect,
	}
}

// Prepare the module
func (r *module) Prepare() []interface{} {
	return []interface{}{
		r.ping,
	}
}

// Destroy dependencies
func (r *module) Destroy() []interface{} {
	return []interface{}{
		r.disconnect,
	}
}
