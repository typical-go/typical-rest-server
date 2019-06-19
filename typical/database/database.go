package database

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/typical-go/typical-rest-server/config"
	"github.com/urfave/cli"

	"github.com/golang-migrate/migrate"

	// load file source driver
	_ "github.com/golang-migrate/migrate/source/file"
)

var DefaultMigrationDirectory = "db/migrate"

// Create database
func Create(conf config.Config) (err error) {
	conn, err := sql.Open("postgres", conf.Postgres.ConnectionStringTemplate1())
	if err != nil {
		return
	}
	defer conn.Close()

	query := fmt.Sprintf(`CREATE DATABASE "%s"`, conf.Postgres.DbName)
	fmt.Println(query)
	_, err = conn.Exec(query)
	return
}

// Drop database
func Drop(conf config.Config) (err error) {
	conn, err := sql.Open("postgres", conf.Postgres.ConnectionStringTemplate1())
	if err != nil {
		return
	}
	defer conn.Close()

	query := fmt.Sprintf(`DROP DATABASE IF EXISTS "%s"`, conf.Postgres.DbName)
	fmt.Println(query)
	_, err = conn.Exec(query)
	return
}

// Migrate database
func Migrate(conf config.Config, args cli.Args) error {
	source := migrationSource(args)
	log.Printf("Migrate database from source '%s'\n", source)

	migration, err := migrate.New(source, conf.Postgres.ConnectionString())
	if err != nil {
		return err
	}
	defer migration.Close()
	return migration.Up()
}

// Rollback database
func Rollback(conf config.Config, args cli.Args) error {
	source := migrationSource(args)
	log.Printf("Migrate database from source '%s'\n", source)

	migration, err := migrate.New(source, conf.Postgres.ConnectionString())
	if err != nil {
		return err
	}
	defer migration.Close()
	return migration.Down()
}

func migrationSource(args cli.Args) string {
	dir := DefaultMigrationDirectory
	if len(args) > 0 {
		dir = args.First()
	}
	return fmt.Sprintf("file://%s", dir)
}
