package typpostgres

import (
	"database/sql"
	"fmt"
	"os"
	"os/exec"
	"strconv"

	"github.com/golang-migrate/migrate"
	log "github.com/sirupsen/logrus"
)

const (
	createDatabaseScriptTemplate = `CREATE DATABASE "%s"`
	dropDatabaseScriptTemplate   = `DROP DATABASE IF EXISTS "%s"`
)

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

func console(config *Config) (err error) {
	os.Setenv("PGPASSWORD", config.Password)
	cmd := exec.Command("psql", "-h", config.Host, "-p", strconv.Itoa(config.Port), "-U", config.User)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stdout
	cmd.Stdin = os.Stdin

	return cmd.Run()
}
