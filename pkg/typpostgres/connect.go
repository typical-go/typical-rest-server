package typpostgres

import (
	"database/sql"
	"fmt"
)

// DB is postgres database handle
type DB struct {
	*sql.DB
}

func (m *Module) connect(cfg Config) (pgDB *DB, err error) {
	var db *sql.DB
	if db, err = sql.Open("postgres", m.dataSource(cfg)); err != nil {
		err = fmt.Errorf("Posgres: Connect: %w", err)
	}
	pgDB = &DB{db}
	return
}

func (*Module) disconnect(db *DB) (err error) {
	if err = db.Close(); err != nil {
		return fmt.Errorf("Postgres: Disconnect: %w", err)
	}
	return
}

func (m *Module) ping(db *DB) (err error) {
	if err = db.Ping(); err != nil {
		return fmt.Errorf("Postgres: Ping: %w", err)
	}
	return
}

func (*Module) dataSource(c Config) string {
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		c.User, c.Password, c.Host, c.Port, c.DBName)
}

func (*Module) adminDataSource(c Config) string {
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		c.User, c.Password, c.Host, c.Port, "template1")
}
