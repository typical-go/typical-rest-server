package typpostgres

import (
	"github.com/typical-go/typical-go/pkg/typbuildtool"
	"github.com/urfave/cli/v2"
)

func cmdResetDB(c *typbuildtool.Context) *cli.Command {
	return &cli.Command{
		Name:   "reset",
		Usage:  "Reset Database",
		Action: c.ActionFunc("PG", resetDB),
	}
}

func resetDB(c *typbuildtool.CliContext) (err error) {

	if err = dropDB(c); err != nil {
		return
	}
	if err = createDB(c); err != nil {
		return
	}
	if err = migrateDB(c); err != nil {
		return
	}
	if err = seedDB(c); err != nil {
		return
	}
	return

}
