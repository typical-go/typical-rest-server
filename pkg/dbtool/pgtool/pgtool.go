package pgtool

import (
	"database/sql"
	"fmt"
	"io"
	"io/ioutil"
	"os"

	"github.com/golang-migrate/migrate"
	"github.com/golang-migrate/migrate/database/postgres"
	"github.com/typical-go/typical-go/pkg/typgo"
	"github.com/typical-go/typical-rest-server/pkg/dbtool"
	"github.com/urfave/cli/v2"

	// load migration file
	_ "github.com/golang-migrate/migrate/source/file"
)

type (
	// PgTool for postgres
	PgTool struct {
		Name         string
		ConfigFn     func() dbtool.Configurer
		DockerName   string
		MigrationSrc string
		SeedSrc      string
		cfg          *dbtool.Config
	}
)

var _ (typgo.Tasker) = (*PgTool)(nil)

// Stdout standard output
var Stdout io.Writer = os.Stdout

// Task for postgres
func (t *PgTool) Task(sys *typgo.BuildSys) *cli.Command {
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

// Cfg ...
func (t *PgTool) Cfg() *dbtool.Config {
	if t.cfg == nil {
		t.cfg = t.ConfigFn().Config()
	}
	return t.cfg
}

// Console interactice for postgres
func (t *PgTool) Console(c *typgo.Context) error {
	cfg := t.Cfg()
	os.Setenv("PGPASSWORD", cfg.DBPass)
	return c.Execute(&typgo.Bash{
		Name: "docker",
		Args: []string{
			"exec", "-it", t.DockerName,
			"psql",
			"-h", cfg.Host,
			"-p", cfg.Port,
			"-U", cfg.DBUser,
		},
		Stdout: os.Stdout,
		Stderr: os.Stderr,
		Stdin:  os.Stdin,
	})
}

// CreateDB create database
func (t *PgTool) CreateDB(c *typgo.Context) error {
	cfg := t.Cfg()
	conn, err := t.createAdminConn()
	if err != nil {
		return err
	}
	defer conn.Close()

	q := fmt.Sprintf("CREATE DATABASE \"%s\"", cfg.DBName)
	fmt.Fprintln(Stdout, "\npg: "+q)
	_, err = conn.ExecContext(c.Ctx(), q)
	return err
}

// DropDB delete database
func (t *PgTool) DropDB(c *typgo.Context) error {
	cfg := t.Cfg()
	conn, err := t.createAdminConn()
	if err != nil {
		return err
	}
	defer conn.Close()

	q := fmt.Sprintf("DROP DATABASE IF EXISTS \"%s\"", cfg.DBName)
	fmt.Fprintln(Stdout, "\npg: "+q)
	_, err = conn.ExecContext(c.Ctx(), q)
	return err
}

// MigrateDB migrate database
func (t *PgTool) MigrateDB(c *typgo.Context) error {
	fmt.Fprintf(Stdout, "\npg: Migrate '%s'\n", t.MigrationSrc)
	migration, err := t.createMigration()
	if err != nil {
		return err
	}
	defer migration.Close()
	return migration.Up()
}

// RollbackDB rollback database
func (t *PgTool) RollbackDB(c *typgo.Context) error {
	fmt.Fprintf(Stdout, "\npg: Rollback '%s'\n", t.MigrationSrc)
	migration, err := t.createMigration()
	if err != nil {
		return err
	}
	defer migration.Close()
	return migration.Down()
}

// SeedDB seed database
func (t *PgTool) SeedDB(c *typgo.Context) error {
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

func (t *PgTool) createMigration() (*migrate.Migrate, error) {
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

func (t *PgTool) createConn() (*sql.DB, error) {
	cfg := t.Cfg()
	return sql.Open("postgres", fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.DBUser, cfg.DBPass, cfg.Host, cfg.Port, cfg.DBName,
	))
}

func (t *PgTool) createAdminConn() (*sql.DB, error) {
	cfg := t.Cfg()
	connStr := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/template1?sslmode=disable",
		cfg.DBUser, cfg.DBPass, cfg.Host, cfg.Port)
	fmt.Println(connStr)
	return sql.Open("postgres", connStr)
}
