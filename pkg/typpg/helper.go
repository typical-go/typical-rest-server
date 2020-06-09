package typpg

import (
	"database/sql"
	"fmt"
)

// Connect to postgres server
func Connect(cfg *Config) (db *sql.DB, err error) {
	if db, err = sql.Open("postgres", Conn(cfg)); err != nil {
		err = fmt.Errorf("postgres: %w", err)
	}

	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("postgres: %w", err)
	}
	return db, nil
}

// Disconnect from postgres server
func Disconnect(db *sql.DB) (err error) {
	if err = db.Close(); err != nil {
		return fmt.Errorf("postgres: %w", err)
	}
	return
}
