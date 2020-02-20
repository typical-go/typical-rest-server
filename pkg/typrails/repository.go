package typrails

import (
	"errors"

	"github.com/typical-go/typical-go/pkg/typcfg"
	"github.com/urfave/cli/v2"
)

func (r *rails) repositoryCmd() *cli.Command {
	return &cli.Command{
		Name:      "repository",
		Aliases:   []string{"repo"},
		Usage:     "Generate Repository from tablename",
		ArgsUsage: "[table] [entity]",
		Before: func(ctx *cli.Context) error {
			return typcfg.LoadEnvFile()
		},
		Action: r.ActionFunc(r.repository),
	}
}

func (r *rails) repository(ctx *cli.Context, f Fetcher) (err error) {
	var (
		table  string
		entity string
		e      *Entity
	)
	if table = ctx.Args().First(); table == "" {
		return errors.New("Missing 'table': check `./typicalw rails repository help` for more detail")
	}
	if entity = ctx.Args().Get(1); entity == "" {
		return errors.New("Missing 'entity': check `./typicalw rails repository help` for more detail")
	}
	if e, err = f.Fetch(r.ModulePackage, table, entity); err != nil {
		return
	}
	if err = generateTransactional(); err != nil {
		return
	}
	if err = generateRepository(e); err != nil {
		return
	}
	return
}
