package typrails

import (
	"errors"

	"github.com/typical-go/typical-go/pkg/typbuildtool"
	"github.com/typical-go/typical-go/pkg/typcfg"
	"github.com/urfave/cli/v2"
)

func repositoryCmd(c *typbuildtool.Context) *cli.Command {
	return &cli.Command{
		Name:      "repository",
		Aliases:   []string{"repo"},
		Usage:     "Generate Repository from tablename",
		ArgsUsage: "[table] [entity]",
		Before: func(ctx *cli.Context) error {
			return typcfg.LoadEnvFile()
		},
		Action: c.ActionFunc(func(cliCtx *cli.Context, f Fetcher) (err error) {
			var (
				table  string
				entity string
				e      *Entity
				ctx    = cliCtx.Context
			)
			if table = cliCtx.Args().First(); table == "" {
				return errors.New("Missing 'table': check `./typicalw rails repository help` for more detail")
			}
			if entity = cliCtx.Args().Get(1); entity == "" {
				return errors.New("Missing 'entity': check `./typicalw rails repository help` for more detail")
			}
			if e, err = f.Fetch(c.ModulePackage, table, entity); err != nil {
				return
			}
			if err = generateTransactional(ctx); err != nil {
				return
			}
			if err = generateRepository(ctx, e); err != nil {
				return
			}
			return
		}),
	}
}
