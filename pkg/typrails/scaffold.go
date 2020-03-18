package typrails

import (
	"errors"

	"github.com/typical-go/typical-go/pkg/typbuildtool"
	"github.com/urfave/cli/v2"
)

func scaffoldCmd(c *typbuildtool.Context) *cli.Command {
	return &cli.Command{
		Name:      "scaffold",
		Aliases:   []string{"s"},
		Usage:     "Generate CRUD API",
		ArgsUsage: "[table] [entity]",
		Action: func(cliCtx *cli.Context) (err error) {
			var (
				table  string
				entity string
				e      *Entity
				ctx    = cliCtx.Context
			)
			if table = cliCtx.Args().First(); table == "" {
				return errors.New("Missing 'table': check `./typicalw rails scaffold help` for more detail")
			}
			if entity = cliCtx.Args().Get(1); entity == "" {
				return errors.New("Missing 'entity': check `./typicalw rails scaffold help` for more detail")
			}
			// if e, err = f.Fetch(c.ProjectPackage, table, entity); err != nil {
			// 	return
			// }
			if err = generateTransactional(ctx); err != nil {
				return
			}
			if err = generateRepository(ctx, e); err != nil {
				return
			}
			if err = generateService(ctx, e); err != nil {
				return
			}
			if err = generateController(ctx, e); err != nil {
				return
			}
			return
		},
	}
}
