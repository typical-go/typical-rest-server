package xpostgres

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/typical-go/typical-go/appx"

	"github.com/golang-migrate/migrate"
)

// DBInfra postgres database infrastructure
type DBInfra struct {
	appx.DBInfra
	config Config
}

// Create database
func (i DBInfra) Create() (err error) {
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
func (i DBInfra) Drop() (err error) {
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
func (i DBInfra) Migrate(source string) error {
	log.Printf("Migrate database from source '%s'\n", source)

	migration, err := migrate.New(source, i.config.ConnectionString())
	if err != nil {
		return err
	}
	defer migration.Close()
	return migration.Up()
}

// Rollback database
func (i DBInfra) Rollback(source string) error {
	log.Printf("Migrate database from source '%s'\n", source)

	migration, err := migrate.New(source, i.config.ConnectionString())
	if err != nil {
		return err
	}
	defer migration.Close()
	return migration.Down()
}
