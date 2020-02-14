package app

import (
	"github.com/typical-go/typical-go/pkg/typapp"
	"github.com/typical-go/typical-go/pkg/typcore"
	"github.com/typical-go/typical-rest-server/app/config"
	"github.com/urfave/cli/v2"
)

// New application [nowire]
func New() *Module {
	return &Module{}
}

// Module of application
type Module struct{}

// EntryPoint of application
func (m Module) EntryPoint() interface{} {
	return func(s server) error {
		s.Middleware()
		s.Route()
		return s.Start()
	}
}

// Configure application
func (m Module) Configure(loader typcore.ConfigLoader) (prefix string, spec, loadFn interface{}) {
	prefix = "APP"
	spec = &config.Config{}
	loadFn = func() (cfg config.Config, err error) {
		err = loader.Load(prefix, &cfg)
		return
	}
	return
}

// AppCommands return comamnd
func (m Module) AppCommands(c *typapp.Context) []*cli.Command {
	return []*cli.Command{
		{
			Name:   "route",
			Usage:  "Print available API Routes",
			Action: c.ActionFunc(taskRouteList),
		},
	}
}
