package app

import (
	"github.com/typical-go/typical-go/pkg/typcfg"
	"github.com/typical-go/typical-go/pkg/typcli"
	"github.com/typical-go/typical-rest-server/app/config"
	"github.com/urfave/cli"
)

// Module of application
func Module() interface{} {
	return module{
		Configuration: typcfg.Configuration{
			Prefix: "APP",
			Spec:   &config.Config{},
		},
	}
}

type module struct {
	typcfg.Configuration
}

func (m module) Action() interface{} {
	return startServer
}

func (m module) Provide() []interface{} {
	return []interface{}{
		func(loader typcfg.Loader) (cfg config.Config, err error) {
			err = loader.Load(m.Configuration, &cfg)
			return
		},
	}
}

func (m module) AppCommands(c *typcli.ContextCli) []cli.Command {
	return []cli.Command{
		{Name: "route", Usage: "Print available API Routes", Action: c.Action(taskRouteList)},
	}
}
