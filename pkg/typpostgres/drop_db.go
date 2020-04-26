package typpostgres

import (
	"database/sql"
	"fmt"

	"github.com/typical-go/typical-go/pkg/typbuildtool"
	"github.com/urfave/cli/v2"
)

func cmdDropDB(c *typbuildtool.Context) *cli.Command {
	return &cli.Command{
		Name:   "drop",
		Usage:  "Drop Database",
		Action: dropDBAction(c),
	}
}

func dropDBAction(c *typbuildtool.Context) cli.ActionFunc {
	return func(cliCtx *cli.Context) (err error) {
		return dropDB(c.BuildContext(cliCtx))
	}
}

func dropDB(c *typbuildtool.BuildContext) (err error) {
	var (
		conn *sql.DB
		cfg  *Config
	)

	if cfg, err = retrieveConfig(c); err != nil {
		return
	}

	if conn, err = sql.Open("postgres", adminDataSource(cfg)); err != nil {
		return
	}
	defer conn.Close()

	query := fmt.Sprintf(`DROP DATABASE IF EXISTS "%s"`, cfg.DBName)
	c.Infof("Postgres: %s", query)
	_, err = conn.ExecContext(c.Cli.Context, query)
	return
}
