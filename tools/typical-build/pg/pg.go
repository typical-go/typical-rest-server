package pg

import (
	"database/sql"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"

	"github.com/golang-migrate/migrate"
	"github.com/golang-migrate/migrate/database/postgres"
	"github.com/typical-go/typical-go/pkg/execkit"
	"github.com/typical-go/typical-go/pkg/typgo"
	"github.com/typical-go/typical-rest-server/internal/app/infra"
	"github.com/urfave/cli/v2"

	// load migration file
	_ "github.com/golang-migrate/migrate/source/file"
)

type (
	// Command for postgres
	Command struct {
		Name         string
		ConfigFn     func() (*infra.PostgresCfg, error)
		DockerName   string
		MigrationSrc string
		SeedSrc      string
		cfg          *infra.PostgresCfg
	}
)

var _ (typgo.Cmd) = (*Command)(nil)

// Stdout standard output
var Stdout io.Writer = os.Stdout

// Command for postgress
func (t *Command) Command(sys *typgo.BuildSys) *cli.Command {
	var err error
	if t.cfg, err = t.ConfigFn(); err != nil {
		log.Fatal(err.Error())
	}

	return &cli.Command{
		Name:  t.Name,
		Usage: t.Name + " utility",
		Subcommands: []*cli.Command{
			{Name: "create", Usage: "Create database", Action: sys.ExecuteFn(t.CreateDB)},
			{Name: "drop", Usage: "Drop database", Action: sys.ExecuteFn(t.DropDB)},
			{Name: "migrate", Usage: "Migrate database", Action: sys.ExecuteFn(t.MigrateDB)},
			{Name: "rollback", Usage: "Rollback database", Action: sys.ExecuteFn(t.RollbackDB)},
			{Name: "seed", Usage: "Seed database", Action: sys.ExecuteFn(t.SeedDB)},
			{Name: "console", Usage: "Postgres console", Action: sys.ExecuteFn(t.Console)},
		},
	}
}

// Console interactice for postgres
func (t *Command) Console(c *typgo.Context) error {
	os.Setenv("PGPASSWORD", t.cfg.DBPass)
	return c.Execute(&execkit.Command{
		Name: "docker",
		Args: []string{
			"exec", "-it", t.DockerName,
			"psql",
			"-h", t.cfg.Host,
			"-p", t.cfg.Port,
			"-U", t.cfg.DBUser,
		},
		Stdout: os.Stdout,
		Stderr: os.Stderr,
		Stdin:  os.Stdin,
	})
}

// CreateDB create database
func (t *Command) CreateDB(c *typgo.Context) error {
	conn, err := t.createAdminConn()
	if err != nil {
		return err
	}
	defer conn.Close()

	q := fmt.Sprintf("CREATE DATABASE \"%s\"", t.cfg.DBName)
	fmt.Fprintln(Stdout, "\npg: "+q)
	_, err = conn.ExecContext(c.Ctx(), q)
	return err
}

// DropDB delete database
func (t *Command) DropDB(c *typgo.Context) error {
	conn, err := t.createAdminConn()
	if err != nil {
		return err
	}
	defer conn.Close()

	q := fmt.Sprintf("DROP DATABASE IF EXISTS \"%s\"", t.cfg.DBName)
	fmt.Fprintln(Stdout, "\npg: "+q)
	_, err = conn.ExecContext(c.Ctx(), q)
	return err
}

// MigrateDB migrate database
func (t *Command) MigrateDB(c *typgo.Context) error {
	fmt.Fprintf(Stdout, "\npg: Migrate '%s'\n", t.MigrationSrc)
	migration, err := t.createMigration()
	if err != nil {
		return err
	}
	defer migration.Close()
	return migration.Up()
}

// RollbackDB rollback database
func (t *Command) RollbackDB(c *typgo.Context) error {
	fmt.Fprintf(Stdout, "\npg: Rollback '%s'\n", t.MigrationSrc)
	migration, err := t.createMigration()
	if err != nil {
		return err
	}
	defer migration.Close()
	return migration.Down()
}

// SeedDB seed database
func (t *Command) SeedDB(c *typgo.Context) error {
	db, err := t.createConn()
	if err != nil {
		return err
	}
	defer db.Close()

	files, _ := ioutil.ReadDir(t.SeedSrc)
	for _, f := range files {
		filename := fmt.Sprintf("%s/%s", t.SeedSrc, f.Name())
		fmt.Printf("\npg: Seed '%s'\n", filename)
		b, _ := ioutil.ReadFile(filename)
		_, err = db.ExecContext(c.Ctx(), string(b))
		if err != nil {
			return err
		}
	}
	return nil
}

func (t *Command) createMigration() (*migrate.Migrate, error) {
	db, err := t.createConn()
	if err != nil {
		return nil, err
	}
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return nil, err
	}
	return migrate.NewWithDatabaseInstance(t.MigrationSrc, "postgres", driver)
}

func (t *Command) createConn() (*sql.DB, error) {
	return sql.Open("postgres", fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		t.cfg.DBUser, t.cfg.DBPass, t.cfg.Host, t.cfg.Port, t.cfg.DBName,
	))
}

func (t *Command) createAdminConn() (*sql.DB, error) {
	return sql.Open("postgres", fmt.Sprintf(
		"postgres://%s:%s@%s:%s/template1?sslmode=disable",
		t.cfg.DBUser, t.cfg.DBPass, t.cfg.Host, t.cfg.Port))
}
