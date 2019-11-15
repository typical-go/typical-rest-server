package app

import (
	"github.com/typical-go/typical-go/pkg/typicli"
	"github.com/typical-go/typical-go/pkg/typictx"
	"github.com/typical-go/typical-go/pkg/typimodule"
	"github.com/typical-go/typical-rest-server/app/config"
	"github.com/urfave/cli"
)

// Module of application
func Module() typictx.AppModule {
	return applicationModule{
		Configuration: typimodule.Configuration{
			Prefix: "APP",
			Spec:   &config.Config{},
		},
	}
}

type applicationModule struct {
	typimodule.Configuration
}

func (m applicationModule) AppCommands(c *typicli.ContextCli) []cli.Command {
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

func (m applicationModule) Run() interface{} {
	return Start
}

func (m applicationModule) Provide() []interface{} {
	return []interface{}{
		m.loadConfig,
	}
}

func (m applicationModule) loadConfig() (cfg *config.Config, err error) {
	err = m.Configuration.Load()
	cfg = m.Configuration.Spec.(*config.Config)
	return
}
