package pgtool

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/golang-migrate/migrate"
	"github.com/golang-migrate/migrate/database/postgres"
	"github.com/typical-go/typical-go/pkg/oskit"
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
		Config       dbtool.Config
		DockerName   string
		MigrationSrc string
		SeedSrc      string
		cfg          *dbtool.Config
	}
)

var _ (typgo.Tasker) = (*PgTool)(nil)

// Task for postgres
func (t *PgTool) Task(sys *typgo.BuildSys) *cli.Command {
	return &cli.Command{
		Name:  t.Name,
		Usage: t.Name + " utility",
		Subcommands: []*cli.Command{
			{Name: "create", Usage: "Create database", Action: sys.ExecuteFn(t.CreateDB)},
			{Name: "drop", Usage: "Drop database", Action: sys.ExecuteFn(t.DropDB)},
			{Name: "migrate", Usage: "Migrate database", Action: sys.ExecuteFn(t.MigrateDB)},
			{Name: "migration", Usage: "Create Migration file", Action: sys.ExecuteFn(t.MigrationFile)},
			{Name: "rollback", Usage: "Rollback database", Action: sys.ExecuteFn(t.RollbackDB)},
			{Name: "seed", Usage: "Seed database", Action: sys.ExecuteFn(t.SeedDB)},
			{Name: "console", Usage: "Postgres console", Action: sys.ExecuteFn(t.Console)},
		},
	}
}

// Console interactice for postgres
func (t *PgTool) Console(c *typgo.Context) error {
	os.Setenv("PGPASSWORD", t.Config.DBPass)
	return c.Execute(&typgo.Bash{
		Name: "docker",
		Args: []string{
			"exec", "-it", t.DockerName,
			"psql",
			"-h", t.Config.Host,
			"-p", t.Config.Port,
			"-U", t.Config.DBUser,
		},
		Stdout: os.Stdout,
		Stderr: os.Stderr,
		Stdin:  os.Stdin,
	})
}

// CreateDB create database
func (t *PgTool) CreateDB(c *typgo.Context) error {
	conn, err := t.createAdminConn()
	if err != nil {
		return err
	}
	defer conn.Close()

	q := fmt.Sprintf("CREATE DATABASE \"%s\"", t.Config.DBName)
	fmt.Fprintf(oskit.Stdout, "\n%s: %s\n", t.Name, q)
	_, err = conn.ExecContext(c.Ctx(), q)
	return err
}

// DropDB delete database
func (t *PgTool) DropDB(c *typgo.Context) error {
	conn, err := t.createAdminConn()
	if err != nil {
		return err
	}
	defer conn.Close()

	q := fmt.Sprintf("DROP DATABASE IF EXISTS \"%s\"", t.Config.DBName)
	fmt.Fprintf(oskit.Stdout, "\n%s: %s\n", t.Name, q)
	_, err = conn.ExecContext(c.Ctx(), q)
	return err
}

// MigrateDB migrate database
func (t *PgTool) MigrateDB(c *typgo.Context) error {
	fmt.Fprintf(oskit.Stdout, "\n%s: Migrate '%s'\n", t.Name, t.MigrationSrc)
	migration, err := t.createMigration()
	if err != nil {
		return err
	}
	defer migration.Close()
	return migration.Up()
}

// RollbackDB rollback database
func (t *PgTool) RollbackDB(c *typgo.Context) error {
	fmt.Fprintf(oskit.Stdout, "\n%s: Rollback '%s'\n", t.Name, t.MigrationSrc)
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
		fmt.Fprintf(oskit.Stdout, "\n%s: Seed '%s'\n", t.Name, filename)
		b, _ := ioutil.ReadFile(filename)
		_, err = db.ExecContext(c.Ctx(), string(b))
		if err != nil {
			return err
		}
	}
	return nil
}

// MigrationFile seed database
func (t *PgTool) MigrationFile(c *typgo.Context) error {
	args := c.Args().Slice()
	if len(args) < 1 {
		args = []string{"migration"}
	}
	for _, arg := range args {
		dbtool.CreateMigrationFile(t.MigrationSrc, arg)
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
	return migrate.NewWithDatabaseInstance("file://"+t.MigrationSrc, "postgres", driver)
}

func (t *PgTool) createConn() (*sql.DB, error) {
	return sql.Open("postgres", fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		t.Config.DBUser, t.Config.DBPass, t.Config.Host, t.Config.Port, t.Config.DBName,
	))
}

func (t *PgTool) createAdminConn() (*sql.DB, error) {
	return sql.Open("postgres", fmt.Sprintf(
		"postgres://%s:%s@%s:%s/template1?sslmode=disable",
		t.Config.DBUser, t.Config.DBPass, t.Config.Host, t.Config.Port))
}
