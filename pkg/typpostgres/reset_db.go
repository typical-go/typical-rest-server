package typpostgres

import (
	"github.com/typical-go/typical-go/pkg/typbuildtool"
	"github.com/urfave/cli/v2"
)

func cmdResetDB(c *typbuildtool.Context) *cli.Command {
	return &cli.Command{
		Name:   "reset",
		Usage:  "Reset Database",
		Action: resetDBAction(c),
	}
}

func resetDBAction(c *typbuildtool.Context) cli.ActionFunc {
	return func(cliCtx *cli.Context) (err error) {
		bc := c.BuildContext(cliCtx)
		if err = dropDB(bc); err != nil {
			return
		}
		if err = createDB(bc); err != nil {
			return
		}
		if err = migrateDB(bc); err != nil {
			return
		}
		if err = seedDB(bc); err != nil {
			return
		}
		return
	}
}
