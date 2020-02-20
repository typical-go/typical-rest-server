package restserver

import (
	"github.com/typical-go/typical-go/pkg/typapp"
	"github.com/typical-go/typical-go/pkg/typcfg"
	"github.com/typical-go/typical-rest-server/restserver/config"
	"github.com/urfave/cli/v2"
)

// Module of application
type Module struct {
	prefix string
}

// New application [nowire]
func New() *Module {
	return &Module{
		prefix: "APP",
	}
}

// WithPrefix return module with new prefix
func (m *Module) WithPrefix(prefix string) *Module {
	m.prefix = prefix
	return m
}

// EntryPoint of application
func (m *Module) EntryPoint() interface{} {
	return func(s server) error {
		s.Middleware()
		s.Route()
		return s.Start()
	}
}

// Configure application
func (m *Module) Configure(loader typcfg.Loader) *typcfg.Detail {
	return &typcfg.Detail{
		Prefix: m.prefix,
		Spec:   &config.Config{},
		Constructor: func() (cfg config.Config, err error) {
			err = loader.Load(m.prefix, &cfg)
			return
		},
	}
}

// AppCommands return comamnd
func (m *Module) AppCommands(c *typapp.Context) []*cli.Command {
	return []*cli.Command{
		{
			Name:   "route",
			Usage:  "Print available API Routes",
			Action: c.ActionFunc(taskRouteList),
		},
	}
}
