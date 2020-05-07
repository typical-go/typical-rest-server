package typpostgres

import (
	"database/sql"
	"fmt"

	_ "github.com/golang-migrate/migrate/database/postgres"
	_ "github.com/golang-migrate/migrate/source/file"
	_ "github.com/lib/pq"
)

// Connect to postgres server
func Connect(cfg *Config) (db *sql.DB, err error) {
	if db, err = sql.Open("postgres", cfg.ConnStr()); err != nil {
		err = fmt.Errorf("Posgres: Connect: %w", err)
	}
	return
}

// Disconnect to postgres server
func Disconnect(db *sql.DB) (err error) {
	if err = db.Close(); err != nil {
		return fmt.Errorf("Postgres: Disconnect: %w", err)
	}
	return
}

// Ping to postgres server
func Ping(db *sql.DB) (err error) {
	if err = db.Ping(); err != nil {
		return fmt.Errorf("Postgres: Ping: %w", err)
	}
	return
}
