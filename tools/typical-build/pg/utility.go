package pg

import (
	"database/sql"
	"fmt"
	"io/ioutil"

	"github.com/golang-migrate/migrate"
	"github.com/golang-migrate/migrate/database/postgres"
	"github.com/typical-go/typical-go/pkg/typgo"
	"github.com/typical-go/typical-rest-server/internal/infra"
)

type (
	utility struct {
		*infra.PostgresCfg
		migrationSrc string
		seedSrc      string
	}
)

// CreateDB create database
func (t *utility) CreateDB(c *typgo.Context) error {
	conn, err := t.createAdminConn()
	if err != nil {
		return err
	}
	defer conn.Close()

	q := fmt.Sprintf(`CREATE DATABASE "%s"`, t.DBName)
	fmt.Fprintln(Stdout, "\npg: "+q)
	_, err = conn.ExecContext(c.Ctx(), q)
	return err
}

// DropDB delete database
func (t *utility) DropDB(c *typgo.Context) error {
	conn, err := t.createAdminConn()
	if err != nil {
		return err
	}
	defer conn.Close()

	q := fmt.Sprintf(`DROP DATABASE IF EXISTS "%s"`, t.DBName)
	fmt.Fprintln(Stdout, "\npg: "+q)
	_, err = conn.ExecContext(c.Ctx(), q)
	return err
}

// MigrateDB migrate database
func (t *utility) MigrateDB(c *typgo.Context) error {
	fmt.Fprintf(Stdout, "\npg: Migrate '%s'\n", t.migrationSrc)
	migration, err := t.createMigration()
	if err != nil {
		return err
	}
	defer migration.Close()
	return migration.Up()
}

// RollbackDB rollback database
func (t *utility) RollbackDB(c *typgo.Context) error {
	fmt.Fprintf(Stdout, "\npg: Rollback '%s'\n", t.migrationSrc)
	migration, err := t.createMigration()
	if err != nil {
		return err
	}
	defer migration.Close()
	return migration.Down()
}

// SeedDB seed database
func (t *utility) SeedDB(c *typgo.Context) error {
	db, err := t.createConn()
	if err != nil {
		return err
	}
	defer db.Close()

	files, _ := ioutil.ReadDir(t.seedSrc)
	for _, f := range files {
		filename := fmt.Sprintf("%s/%s", t.seedSrc, f.Name())
		fmt.Printf("\npg: Seed '%s'\n", filename)
		b, _ := ioutil.ReadFile(filename)
		_, err = db.ExecContext(c.Ctx(), string(b))
		if err != nil {
			return err
		}
	}
	return nil
}

func (t *utility) createMigration() (*migrate.Migrate, error) {
	db, err := t.createConn()
	if err != nil {
		return nil, err
	}
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return nil, err
	}
	return migrate.NewWithDatabaseInstance(t.migrationSrc, "postgres", driver)
}

func (t *utility) createConn() (*sql.DB, error) {
	return sql.Open("postgres", fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		t.DBUser, t.DBPass, t.Host, t.Port, t.DBName,
	))
}

func (t *utility) createAdminConn() (*sql.DB, error) {
	return sql.Open("postgres", fmt.Sprintf(
		"postgres://%s:%s@%s:%s/template1?sslmode=disable",
		t.DBUser, t.DBPass, t.Host, t.Port))
}
