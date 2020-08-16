package pg

import (
	"database/sql"
	"fmt"
	"io"
	"io/ioutil"
	"os"

	"github.com/golang-migrate/migrate"
	"github.com/golang-migrate/migrate/database/postgres"
	"github.com/typical-go/typical-go/pkg/typgo"
	"github.com/typical-go/typical-rest-server/internal/infra"

	// load migration file
	_ "github.com/golang-migrate/migrate/source/file"
)

// Stdout standard output
var Stdout io.Writer = os.Stdout

type (
	// Tool postgres
	Tool struct {
		*infra.PostgresCfg
		migrationSrc string
		seedSrc      string
	}
)

// CreateDB create database
func (t *Tool) CreateDB(c *typgo.Context) error {
	conn, err := t.createAdminConn()
	if err != nil {
		return err
	}
	defer conn.Close()

	q := fmt.Sprintf(`CREATE DATABASE "%s"`, t.DBName)
	fmt.Fprintln(Stdout, q)
	_, err = conn.ExecContext(c.Ctx(), q)
	return err
}

// DropDB delete database
func (t *Tool) DropDB(c *typgo.Context) error {
	conn, err := t.createAdminConn()
	if err != nil {
		return err
	}
	defer conn.Close()

	q := fmt.Sprintf(`DROP DATABASE IF EXISTS "%s"`, t.DBName)
	fmt.Fprintln(Stdout, q)
	_, err = conn.ExecContext(c.Ctx(), q)
	return err
}

// MigrateDB migrate database
func (t *Tool) MigrateDB(c *typgo.Context) error {
	fmt.Fprintf(Stdout, "Migrate '%s'\n", t.migrationSrc)
	migration, err := t.createMigration()
	if err != nil {
		return err
	}
	defer migration.Close()
	return migration.Up()
}

// RollbackDB rollback database
func (t *Tool) RollbackDB(c *typgo.Context) error {
	migration, err := t.createMigration()
	if err != nil {
		return err
	}
	defer migration.Close()
	return migration.Down()
}

// SeedDB seed database
func (t *Tool) SeedDB(c *typgo.Context) error {
	db, err := t.createConn()
	if err != nil {
		return err
	}
	defer db.Close()

	fmt.Printf("Seed '%s'\n", t.seedSrc)
	files, _ := ioutil.ReadDir(t.seedSrc)
	for _, f := range files {
		b, err := ioutil.ReadFile(t.seedSrc + "/" + f.Name())
		if err != nil {
			return err
		}
		_, err = db.ExecContext(c.Ctx(), string(b))
		if err != nil {
			return err
		}
	}
	return nil
}

func (t *Tool) createMigration() (*migrate.Migrate, error) {
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

func (t *Tool) createConn() (*sql.DB, error) {
	return sql.Open("postgres", fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		t.User, t.Pass, t.Host, t.Port, t.DBName,
	))
}

func (t *Tool) createAdminConn() (*sql.DB, error) {
	return sql.Open("postgres", fmt.Sprintf(
		"postgres://%s:%s@%s:%s/template1?sslmode=disable",
		t.User, t.Pass, t.Host, t.Port))
}
