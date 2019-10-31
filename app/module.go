package app

import (
	"github.com/kelseyhightower/envconfig"
	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/typiobj"
	"github.com/typical-go/typical-rest-server/app/config"
	"github.com/urfave/cli"
	"go.uber.org/dig"
)

// Module of application
func Module() interface{} {
	return applicationModule{
		Configuration: typiobj.Configuration{
			Prefix: "APP",
			Spec:   &config.Config{},
		},
	}
}

type applicationModule struct {
	typiobj.Configuration
}

func (m applicationModule) CommandLine() cli.Command {
	return cli.Command{
		Name: "route", Description: "Print available API Routes", Action: Route,
	}
}

func (m applicationModule) Prepare() []interface{} {
	return []interface{}{
		Routes,
		Middlewares,
	}
}

func (m applicationModule) Run(c *dig.Container) (err error) {
	return c.Invoke(Start)
}

func (m applicationModule) Provide() []interface{} {
	return []interface{}{
		m.loadConfig,
	}
}

func (m applicationModule) loadConfig() (cfg *config.Config, err error) {
	cfg = new(config.Config)
	err = envconfig.Process(m.Configure().Prefix, cfg)
	return
}
