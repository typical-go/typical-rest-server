package app

import (
	"github.com/typical-go/typical-go/pkg/typcfg"
	"github.com/typical-go/typical-go/pkg/typcli"
	"github.com/typical-go/typical-rest-server/app/config"
	"github.com/urfave/cli"
)

// Module of application
func Module() interface{} {
	return applicationModule{
		Configuration: typcfg.Configuration{
			Prefix: "APP",
			Spec:   &config.Config{},
		},
	}
}

type applicationModule struct {
	typcfg.Configuration
}

func (m applicationModule) AppCommands(c *typcli.ContextCli) []cli.Command {
	return []cli.Command{
		{Name: "route", Usage: "Print available API Routes", Action: c.Action(taskRouteList)},
	}
}

func (m applicationModule) Prepare() []interface{} {
	return []interface{}{
		Routes,
		Middlewares,
	}
}

func (m applicationModule) Action() interface{} {
	return Start
}

func (m applicationModule) Provide() []interface{} {
	return []interface{}{
		m.loadConfig,
	}
}

func (m applicationModule) loadConfig(loader typcfg.Loader) (cfg config.Config, err error) {
	err = loader.Load(m.Configuration, &cfg)
	return
}
