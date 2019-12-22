package typrails

import (
	"github.com/typical-go/typical-go/pkg/utility/common"
	"github.com/urfave/cli/v2"
)

func (r *rails) scaffoldCmd() *cli.Command {
	return &cli.Command{
		Name:      "scaffold",
		Aliases:   []string{"s"},
		Usage:     "Generate CRUD API",
		ArgsUsage: "[table name]",
		Before: func(ctx *cli.Context) error {
			return common.LoadEnvFile()
		},
		Action: r.PreparedAction(r.scaffold),
	}
}

func (r *rails) scaffold(ctx *cli.Context, f Fetcher) (err error) {
	tableName := ctx.Args().First()
	if tableName == "" {
		return cli.ShowCommandHelp(ctx, "rails")
	}
	var e *Entity
	if e, err = f.Fetch(r.Package, tableName); err != nil {
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
