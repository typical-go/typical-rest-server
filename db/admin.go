package db

import (
	"database/sql"
	"fmt"

	"github.com/imantung/go-helper/dbkit"
	"github.com/imantung/typical-go-server/config"

	"github.com/golang-migrate/migrate"
	"github.com/golang-migrate/migrate/database/postgres"
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
func Migrate(conn *sql.DB) error {
	migration, err := newMigration(conn, "db/migrate")
	if err != nil {
		return err
	}
	defer migration.Close()

	return migration.Up()
}

// Rollback database
func Rollback(conn *sql.DB) error {
	migration, err := newMigration(conn, "db/migrate")
	if err != nil {
		return err
	}
	defer migration.Close()
	return migration.Down()
}

func executeFromTemplateDB(conf config.Config, query string) (err error) {
	pgConf := dbkit.PgConfig{
		Host:     conf.DbHost,
		Port:     conf.DbPort,
		DbName:   "template1",
		User:     conf.DbUser,
		Password: conf.DbPassword,
		SslMode:  "disable",
	}
	conn, err := sql.Open("postgres", pgConf.ConnectionString())
	if err != nil {
		return
	}
	_, err = conn.Exec(query)
	return
}

func newMigration(conn *sql.DB, dir string) (m *migrate.Migrate, err error) {
	sourceURL := fmt.Sprintf("file://%s", dir)
	driver, err := postgres.WithInstance(conn, &postgres.Config{})
	if err != nil {
		return
	}
	m, err = migrate.NewWithDatabaseInstance(sourceURL, "postgres", driver)
	return
}
