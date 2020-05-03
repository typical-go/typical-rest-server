package typpostgres

import (
	"database/sql"
	"fmt"

	"github.com/typical-go/typical-go/pkg/typbuildtool"
	"github.com/urfave/cli/v2"
)

func cmdCreateDB(c *typbuildtool.Context) *cli.Command {
	return &cli.Command{
		Name:   "create",
		Usage:  "Create New Database",
		Action: c.ActionFunc("PG", createDB),
	}
}

func createDB(c *typbuildtool.CliContext) (err error) {
	var (
		conn *sql.DB
		cfg  *Config
	)

	if cfg, err = retrieveConfig(); err != nil {
		return
	}

	if conn, err = sql.Open("postgres", adminDataSource(cfg)); err != nil {
		return
	}
	defer conn.Close()

	if err = conn.PingContext(c.Context); err != nil {
		return
	}

	query := fmt.Sprintf(`CREATE DATABASE "%s"`, cfg.DBName)
	c.Infof("Postgres: %s", query)
	_, err = conn.ExecContext(c.Context, query)
	return
}
