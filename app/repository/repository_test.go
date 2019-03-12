package repository

import (
	"database/sql"

	"github.com/golang-migrate/migrate"
	"github.com/golang-migrate/migrate/database/sqlite3"
)

const (
	migrationSource = "file://../../db/migrate"
)

// MigrateTestDB to in-memory database
func migrateTestDB(source string) (m *migrate.Migrate, conn *sql.DB, err error) {
	conn, _ = sql.Open("sqlite3", ":memory:")
	driver, _ := sqlite3.WithInstance(conn, &sqlite3.Config{})
	m, err = migrate.NewWithDatabaseInstance(source, "sqlite3", driver)
	if err != nil {
		return
	}

	err = m.Up()
	return
}
