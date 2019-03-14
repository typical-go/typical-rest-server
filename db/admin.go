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
func Create(conf config.Config) (err error) {
	conn, err := sql.Open("postgres", connectionStringWithDBName(conf, "template1"))
	if err != nil {
		return
	}
	defer conn.Close()

	query := fmt.Sprintf(`CREATE DATABASE "%s"`, conf.DbName)
	fmt.Println(query)
	_, err = conn.Exec(query)
	return
}

// Drop database
func Drop(conf config.Config) (err error) {
	conn, err := sql.Open("postgres", connectionStringWithDBName(conf, "template1"))
	if err != nil {
		return
	}
	defer conn.Close()

	query := fmt.Sprintf(`DROP DATABASE IF EXISTS "%s"`, conf.DbName)
	fmt.Println(query)
	_, err = conn.Exec(query)
	return
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

// ResetTestDB
func ResetTestDB(conf config.Config, source string) (err error) {
	conn, err := sql.Open("postgres", connectionStringWithDBName(conf, "template1"))
	if err != nil {
		return
	}
	defer conn.Close()

	_, err = conn.Exec(fmt.Sprintf(`DROP DATABASE IF EXISTS "%s"`, conf.DbName))
	if err != nil {
		return
	}
	_, err = conn.Exec(fmt.Sprintf(`CREATE DATABASE "%s"`, conf.DbName))
	if err != nil {
		return
	}

	migration, err := migrate.New(source, connectionString(conf))
	if err != nil {
		return err
	}
	defer migration.Close()
	return migration.Up()
}

func migrationSource(args cli.Args) string {
	dir := config.DefaultMigrationDirectory
	if len(args) > 0 {
		dir = args.First()
	}
	return fmt.Sprintf("file://%s", dir)
}
