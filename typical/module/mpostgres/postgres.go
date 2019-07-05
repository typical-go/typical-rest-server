package mpostgres

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/golang-migrate/migrate"
)

func createDB(config Config) (err error) {
	conn, err := sql.Open("postgres", config.ConnectionStringNoDB())
	if err != nil {
		return
	}
	defer conn.Close()

	query := fmt.Sprintf(`CREATE DATABASE "%s"`, config.DbName)
	fmt.Println(query)
	_, err = conn.Exec(query)
	return
}

func dropDB(config Config) (err error) {
	conn, err := sql.Open("postgres", config.ConnectionStringNoDB())
	if err != nil {
		return
	}
	defer conn.Close()

	query := fmt.Sprintf(`DROP DATABASE IF EXISTS "%s"`, config.DbName)
	fmt.Println(query)
	_, err = conn.Exec(query)
	return
}

func migrateDB(config Config) error {
	source := config.MigrationSource
	log.Printf("Migrate database from source '%s'\n", source)

	migration, err := migrate.New(source, config.ConnectionString())
	if err != nil {
		return err
	}
	defer migration.Close()
	return migration.Up()
}

func rollbackDB(config Config) error {
	source := config.MigrationSource
	log.Printf("Migrate database from source '%s'\n", source)

	migration, err := migrate.New(source, config.ConnectionString())
	if err != nil {
		return err
	}
	defer migration.Close()
	return migration.Down()
}
