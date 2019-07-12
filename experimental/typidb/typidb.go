package typidb

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/golang-migrate/migrate"
)

// Config for database configuration
type Config interface {
	DatabaseName() string
	DataSource() string
	AdminDataSource() string
	DriverName() string
	MigrationSource() string
}

// Tool contain tool of database operation
type Tool struct {
	CreateDatabaseScriptTemplate string
	DropDatabaseScriptTemplate   string
}

// NewPostgresTool return new instance of Tool for Postgres
func NewPostgresTool() *Tool {
	return &Tool{
		CreateDatabaseScriptTemplate: `CREATE DATABASE "%s"`,
		DropDatabaseScriptTemplate:   `DROP DATABASE IF EXISTS "%s"`,
	}
}

// CreateDB is tool to create new database
func (t *Tool) CreateDB(config Config) (err error) {
	query := fmt.Sprintf(t.CreateDatabaseScriptTemplate, config.DatabaseName())
	log.Printf(query)

	conn, err := sql.Open(config.DriverName(), config.AdminDataSource())
	if err != nil {
		return
	}
	defer conn.Close()

	_, err = conn.Exec(query)
	return
}

// DropDB is tool to drop database
func (t *Tool) DropDB(config Config) (err error) {
	query := fmt.Sprintf(t.DropDatabaseScriptTemplate, config.DatabaseName())
	log.Printf(query)

	conn, err := sql.Open(config.DriverName(), config.AdminDataSource())
	if err != nil {
		return
	}
	defer conn.Close()

	_, err = conn.Exec(query)
	return
}

// MigrateDB is tool to migrate database
func (t *Tool) MigrateDB(config Config) error {
	source := config.MigrationSource()
	log.Printf("Migrate database from source '%s'\n", source)

	migration, err := migrate.New(source, config.DataSource())
	if err != nil {
		return err
	}
	defer migration.Close()
	return migration.Up()
}

// RollbackDB is tool to rollback database
func (t *Tool) RollbackDB(config Config) error {
	source := config.MigrationSource()
	log.Printf("Migrate database from source '%s'\n", source)

	migration, err := migrate.New(source, config.DataSource())
	if err != nil {
		return err
	}
	defer migration.Close()
	return migration.Down()
}
