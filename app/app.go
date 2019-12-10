package app

import (
	"github.com/typical-go/typical-go/pkg/typobj"
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

func (m module) AppCommands(c typobj.Cli) []*cli.Command {
	return []*cli.Command{
		{Name: "route", Usage: "Print available API Routes", Action: c.PreparedAction(taskRouteList)},
	}
}

func (m module) Configure() (prefix string, spec, loadFn interface{}) {
	prefix = "APP"
	spec = &config.Config{}
	loadFn = func(loader typobj.Loader) (cfg config.Config, err error) {
		err = loader.Load(prefix, &cfg)
		return
	}
	return
}
