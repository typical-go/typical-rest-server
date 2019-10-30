package app

import (
	"github.com/kelseyhightower/envconfig"
	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/typictx"
	"github.com/typical-go/typical-rest-server/app/config"
	"github.com/urfave/cli"
	"go.uber.org/dig"
)

// Module of application
func Module() interface{} {
	return applicationModule{
		Configuration: typictx.Configuration{
			Prefix: "APP",
			Spec:   &config.Config{},
		},
	}
}

type applicationModule struct {
	typictx.Configuration
}

func (m applicationModule) CommandLine() cli.Command {
	return cli.Command{
		Name: "route", Description: "Print available API Routes", Action: Route,
	}
}

func (m applicationModule) Prepare(c *dig.Container) (err error) {
	if err = c.Invoke(Middlewares); err != nil {
		return
	}
	if err = c.Invoke(Routes); err != nil {
		return
	}
	return
}

func (m applicationModule) Run(c *dig.Container) (err error) {
	return c.Invoke(Start)
}

func (m applicationModule) Construct(c *dig.Container) (err error) {
	return c.Provide(m.loadConfig)
}

func (m applicationModule) loadConfig() (cfg *config.Config, err error) {
	cfg = new(config.Config)
	err = envconfig.Process(m.Configure().Prefix, cfg)
	return
}
