package typrails

import (
	"errors"

	"github.com/typical-go/typical-go/pkg/typcfg"
	"github.com/urfave/cli/v2"
)

func (r *rails) scaffoldCmd() *cli.Command {
	return &cli.Command{
		Name:      "scaffold",
		Aliases:   []string{"s"},
		Usage:     "Generate CRUD API",
		ArgsUsage: "[table] [entity]",
		Before: func(ctx *cli.Context) error {
			return typcfg.LoadEnvFile()
		},
		Action: r.ActionFunc(r.scaffold),
	}
}

func (r *rails) scaffold(ctx *cli.Context, f Fetcher) (err error) {
	var (
		table  string
		entity string
		e      *Entity
	)
	if table = ctx.Args().First(); table == "" {
		return errors.New("Missing 'table': check `./typicalw rails scaffold help` for more detail")
	}
	if entity = ctx.Args().Get(1); entity == "" {
		return errors.New("Missing 'entity': check `./typicalw rails scaffold help` for more detail")
	}
	if e, err = f.Fetch(r.Package, table, entity); err != nil {
		return
	}
	if err = generateTransactional(); err != nil {
		return
	}
	if err = generateRepository(e); err != nil {
		return
	}
	if err = generateService(e); err != nil {
		return
	}
	if err = generateController(e); err != nil {
		return
	}
	return
}
