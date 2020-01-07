package typpostgres

import (
	"github.com/typical-go/typical-go/pkg/typcore"
	"github.com/urfave/cli/v2"
)

func (m *Module) resetCmd(c *typcore.BuildContext) *cli.Command {
	return &cli.Command{
		Name:   "reset",
		Usage:  "Reset Database",
		Action: c.ActionFunc(m.reset),
	}
}

func (m *Module) reset(cfg Config) (err error) {
	if err = m.drop(cfg); err != nil {
		return
	}
	if err = m.create(cfg); err != nil {
		return
	}
	if err = m.migrate(cfg); err != nil {
		return
	}
	if err = m.seed(cfg); err != nil {
		return
	}
	return
}
