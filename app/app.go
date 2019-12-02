package app

import (
	"github.com/typical-go/typical-go/pkg/typcfg"
	"github.com/typical-go/typical-go/pkg/typcli"
	"github.com/typical-go/typical-rest-server/app/config"
	"github.com/urfave/cli/v2"
)

// Module of application
func Module() interface{} {
	return module{}
}

type module struct{}

func (m module) Action() interface{} {
	return startServer
}

func (m module) Commands(c *typcli.AppCli) []*cli.Command {
	return []*cli.Command{
		{Name: "route", Usage: "Print available API Routes", Action: c.Action(taskRouteList)},
	}
}

func (m module) Configure() (prefix string, spec, loadFn interface{}) {
	prefix = "APP"
	spec = &config.Config{}
	loadFn = func(loader typcfg.Loader) (cfg config.Config, err error) {
		err = loader.Load(prefix, &cfg)
		return
	}
	return
}
