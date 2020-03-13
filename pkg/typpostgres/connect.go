package typpostgres

import (
	"database/sql"
	"fmt"
)

func (m *Postgres) connect(cfg *Config) (pgDB *DB, err error) {
	var db *sql.DB
	if db, err = sql.Open("postgres", m.dataSource(cfg)); err != nil {
		err = fmt.Errorf("Posgres: Connect: %w", err)
	}
	pgDB = NewDB(db)
	return
}

func (*Postgres) disconnect(db *DB) (err error) {
	if err = db.Close(); err != nil {
		return fmt.Errorf("Postgres: Disconnect: %w", err)
	}
	return
}

func (m *Postgres) ping(db *DB) (err error) {
	if err = db.Ping(); err != nil {
		return fmt.Errorf("Postgres: Ping: %w", err)
	}
	return
}

func (*Postgres) dataSource(c *Config) string {
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		c.User, c.Password, c.Host, c.Port, c.DBName)
}

func (*Postgres) adminDataSource(c *Config) string {
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		c.User, c.Password, c.Host, c.Port, "template1")
}
