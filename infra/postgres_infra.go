package infra

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/golang-migrate/migrate"

	"github.com/urfave/cli"
)

// PostgresInfra postgres database infrastructure
type PostgresInfra struct {
	config PostgresConfig
}

// DefaultMigrationDirectory default migration directory
var DefaultMigrationDirectory = "db/migration"

// Create database
func (i PostgresInfra) Create() (err error) {
	conn, err := sql.Open("postgres", i.config.ConnectionStringNoDB())
	if err != nil {
		return
	}
	defer conn.Close()

	query := fmt.Sprintf(`CREATE DATABASE "%s"`, i.config.DbName)
	fmt.Println(query)
	_, err = conn.Exec(query)
	return
}

// Drop database
func (i PostgresInfra) Drop() (err error) {
	conn, err := sql.Open("postgres", i.config.ConnectionStringNoDB())
	if err != nil {
		return
	}
	defer conn.Close()

	query := fmt.Sprintf(`DROP DATABASE IF EXISTS "%s"`, i.config.DbName)
	fmt.Println(query)
	_, err = conn.Exec(query)
	return
}

// Migrate database
func (i PostgresInfra) Migrate(args cli.Args) error {
	source := i.migrationSource(args)
	log.Printf("Migrate database from source '%s'\n", source)

	migration, err := migrate.New(source, i.config.ConnectionString())
	if err != nil {
		return err
	}
	defer migration.Close()
	return migration.Up()
}

// Rollback database
func (i PostgresInfra) Rollback(args cli.Args) error {
	source := i.migrationSource(args)
	log.Printf("Migrate database from source '%s'\n", source)

	migration, err := migrate.New(source, i.config.ConnectionString())
	if err != nil {
		return err
	}
	defer migration.Close()
	return migration.Down()
}

func (i PostgresInfra) migrationSource(args cli.Args) string {
	dir := DefaultMigrationDirectory
	if len(args) > 0 {
		dir = args.First()
	}
	return fmt.Sprintf("file://%s", dir)
}
