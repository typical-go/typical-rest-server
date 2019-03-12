package db

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/imantung/typical-go-server/config"
	"github.com/urfave/cli"

	"github.com/golang-migrate/migrate"
	_ "github.com/golang-migrate/migrate/database/postgres"
	_ "github.com/golang-migrate/migrate/source/file"
)

// Create database
func Create(conf config.Config) error {
	query := fmt.Sprintf(`CREATE DATABASE "%s"`, conf.DbName)
	fmt.Println(query)
	return executeFromTemplateDB(conf, query)
}

// Drop database
func Drop(conf config.Config) error {
	query := fmt.Sprintf(`DROP DATABASE "%s"`, conf.DbName)
	fmt.Println(query)
	return executeFromTemplateDB(conf, query)
}

// Migrate database
func Migrate(conf config.Config, args cli.Args) error {
	source := migrationSource(args)
	log.Printf("Migrate database from source '%s'\n", source)

	migration, err := migrate.New(source, connectionString(conf))
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

	migration, err := migrate.New(source, connectionString(conf))
	if err != nil {
		return err
	}
	defer migration.Close()
	return migration.Down()
}

func executeFromTemplateDB(conf config.Config, query string) (err error) {
	conn, err := sql.Open("postgres", connectionStringWithDBName(conf, "template1"))
	if err != nil {
		return
	}
	_, err = conn.Exec(query)
	return
}

func migrationSource(args cli.Args) string {
	dir := config.DefaultMigrationDirectory
	if len(args) > 0 {
		dir = args.First()
	}
	return fmt.Sprintf("file://%s", dir)
}
