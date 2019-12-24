package typrails

import (
	"github.com/typical-go/typical-go/pkg/common"
	"github.com/urfave/cli/v2"
)

func (r *rails) repositoryCmd() *cli.Command {
	return &cli.Command{
		Name:      "repository",
		Aliases:   []string{"repo"},
		Usage:     "Generate Repository from tablename",
		ArgsUsage: "[table name]",
		Before: func(ctx *cli.Context) error {
			return common.LoadEnvFile()
		},
		Action: r.PreparedAction(r.repository),
	}
}

func (r *rails) repository(ctx *cli.Context, f Fetcher) (err error) {
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
	return
}
