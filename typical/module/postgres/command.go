package postgres

import (
	"database/sql"
	"fmt"

	"github.com/golang-migrate/migrate"
	log "github.com/sirupsen/logrus"
	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/typictx"
)

const (
	createDatabaseScriptTemplate = `CREATE DATABASE "%s"`
	dropDatabaseScriptTemplate   = `DROP DATABASE IF EXISTS "%s"`
)

// CreateDB is tool to create new database
func CreateDB(ctx *typictx.ActionContext) (err error) {
	return ctx.Container().Invoke(createDB)
}

func createDB(config *Config) (err error) {
	query := fmt.Sprintf(createDatabaseScriptTemplate, config.DatabaseName())
	log.Infof("Postgres: %s", query)

	conn, err := sql.Open(config.DriverName(), config.AdminDataSource())
	if err != nil {
		return
	}
	defer conn.Close()

	_, err = conn.Exec(query)
	return
}

// DropDB is tool to drop database
func DropDB(ctx *typictx.ActionContext) (err error) {
	return ctx.Container().Invoke(dropDB)
}

func dropDB(config *Config) (err error) {
	query := fmt.Sprintf(dropDatabaseScriptTemplate, config.DatabaseName())
	log.Infof("Postgres: %s", query)

	conn, err := sql.Open(config.DriverName(), config.AdminDataSource())
	if err != nil {
		return
	}
	defer conn.Close()

	_, err = conn.Exec(query)
	return
}

// MigrateDB is tool to migrate database
func MigrateDB(ctx *typictx.ActionContext) (err error) {
	return ctx.Container().Invoke(migrateDB)
}

func migrateDB(config *Config) error {
	sourceURL := fmt.Sprintf("file://%s", config.MigrationSource())
	log.Infof("Migrate database from source '%s'\n", sourceURL)

	migration, err := migrate.New(sourceURL, config.DataSource())
	if err != nil {
		return err
	}
	defer migration.Close()
	return migration.Up()
}

// RollbackDB is tool to rollback database
func RollbackDB(ctx *typictx.ActionContext) (err error) {
	return ctx.Container().Invoke(rollbackDB)
}

func rollbackDB(config *Config) error {
	sourceURL := fmt.Sprintf("file://%s", config.MigrationSource())
	log.Infof("Migrate database from source '%s'\n", sourceURL)

	migration, err := migrate.New(sourceURL, config.DataSource())
	if err != nil {
		return err
	}
	defer migration.Close()
	return migration.Down()
}
