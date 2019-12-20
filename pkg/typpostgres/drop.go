package typpostgres

import (
	"database/sql"
	"fmt"

	log "github.com/sirupsen/logrus"
	"github.com/typical-go/typical-go/pkg/typcore"
	"github.com/urfave/cli/v2"
)

func (m module) dropCmd(c *typcore.Context) *cli.Command {
	return &cli.Command{
		Name:   "drop",
		Usage:  "Drop Database",
		Action: c.Action(m, m.drop),
	}
}

func (m module) drop(cfg Config) (err error) {
	var conn *sql.DB
	query := fmt.Sprintf(`DROP DATABASE IF EXISTS "%s"`, cfg.DBName)
	log.Infof("Postgres: %s", query)
	if conn, err = sql.Open("postgres", m.adminDataSource(cfg)); err != nil {
		return
	}
	defer conn.Close()
	_, err = conn.Exec(query)
	return
}
