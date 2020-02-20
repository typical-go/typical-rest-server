package typpostgres

import (
	"database/sql"
	"fmt"

	log "github.com/sirupsen/logrus"
)

func (m *Postgres) create(cfg Config) (err error) {
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
