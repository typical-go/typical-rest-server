package typpostgres

import (
	"database/sql"
	"fmt"

	"github.com/typical-go/typical-go/pkg/typbuildtool"
	"github.com/urfave/cli/v2"
)

func cmdCreateDB(c *typbuildtool.Context) *cli.Command {
	return &cli.Command{
		Name:  "create",
		Usage: "Create New Database",
		Action: func(cliCtx *cli.Context) (err error) {
			return createDB(c.BuildContext(cliCtx))
		},
	}
}

func createDB(c *typbuildtool.BuildContext) (err error) {
	var conn *sql.DB
	var cfg *Config

	if cfg, err = retrieveConfig(c); err != nil {
		return
	}

	if conn, err = sql.Open("postgres", adminDataSource(cfg)); err != nil {
		return
	}
	defer conn.Close()
	if err = conn.Ping(); err != nil {
		return
	}

	query := fmt.Sprintf(`CREATE DATABASE "%s"`, cfg.DBName)
	c.Infof("Postgres: %s", query)
	_, err = conn.Exec(query)
	return
}
