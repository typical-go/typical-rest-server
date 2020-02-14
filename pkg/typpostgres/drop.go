package typpostgres

import (
	"database/sql"
	"fmt"

	log "github.com/sirupsen/logrus"
)

func (m *Module) drop(cfg Config) (err error) {
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
