package typpostgres

import (
	"database/sql"
	"fmt"

	log "github.com/sirupsen/logrus"
	"github.com/typical-go/typical-go/pkg/typcore"
	"github.com/urfave/cli/v2"
)

func (m Module) createCmd(c *typcore.Context) *cli.Command {
	return &cli.Command{
		Name:   "create",
		Usage:  "Create New Database",
		Action: c.Action(m, m.create),
	}
}

func (m Module) create(cfg Config) (err error) {
	var conn *sql.DB
	query := fmt.Sprintf(`CREATE DATABASE "%s"`, cfg.DBName)
	log.Infof("Postgres: %s", query)
	if conn, err = sql.Open("postgres", m.adminDataSource(cfg)); err != nil {
		return
	}
	defer conn.Close()
	if err = conn.Ping(); err != nil {
		return
	}
	_, err = conn.Exec(query)
	return
}
