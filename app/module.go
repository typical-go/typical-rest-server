package app

import (
	"github.com/typical-go/typical-go/pkg/typcore"
	"github.com/typical-go/typical-rest-server/app/config"
	"github.com/urfave/cli/v2"
)

// New application
func New() interface{} {
	return module{}
}

type module struct{}

func (m module) Action() interface{} {
	return func(s server) error {
		s.Middleware()
		s.Route()
		return s.Start()
	}
}

func (m module) Configure() (prefix string, spec, loadFn interface{}) {
	prefix = "APP"
	spec = &config.Config{}
	loadFn = func(loader typcore.ConfigLoader) (cfg config.Config, err error) {
		err = loader.Load(prefix, &cfg)
		return
	}
	return
}

func (m module) AppCommands(c *typcore.Context) []*cli.Command {
	return []*cli.Command{
		{Name: "route", Usage: "Print available API Routes", Action: c.PreparedAction(taskRouteList)},
	}
}
