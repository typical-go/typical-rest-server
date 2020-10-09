package mysqltool

import (
	"database/sql"
	"fmt"
	"io"
	"io/ioutil"
	"os"

	"github.com/golang-migrate/migrate"
	"github.com/golang-migrate/migrate/database/mysql"
	"github.com/typical-go/typical-go/pkg/execkit"
	"github.com/typical-go/typical-go/pkg/typgo"
	"github.com/typical-go/typical-rest-server/pkg/dbtool"
	"github.com/urfave/cli/v2"
)

type (
	// MySQLTool for postgres
	MySQLTool struct {
		Name         string
		ConfigFn     func() dbtool.Configurer
		DockerName   string
		MigrationSrc string
		SeedSrc      string
		cfg          *dbtool.Config
	}
)

var _ (typgo.Cmd) = (*MySQLTool)(nil)

// Stdout standard output
var Stdout io.Writer = os.Stdout

// Command for postgress
func (t *MySQLTool) Command(sys *typgo.BuildSys) *cli.Command {

	t.cfg = t.ConfigFn().Config()
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
func (t *MySQLTool) Console(c *typgo.Context) error {
	os.Setenv("PGPASSWORD", t.cfg.DBPass)
	return c.Execute(&execkit.Command{
		Name: "docker",
		Args: []string{
			"exec", "-it", t.DockerName,
			"mysql",
			"-h", t.cfg.Host, // host
			"-P", t.cfg.Port, // port
			"-u", t.cfg.DBUser, // user
			fmt.Sprintf("-p%s", t.cfg.DBPass), // password flag can't be spaced
		},
		Stdout: os.Stdout,
		Stderr: os.Stderr,
		Stdin:  os.Stdin,
	})
}

// CreateDB create database
func (t *MySQLTool) CreateDB(c *typgo.Context) error {
	conn, err := t.createAdminConn()
	if err != nil {
		return err
	}
	defer conn.Close()

	q := fmt.Sprintf("CREATE DATABASE `%s`", t.cfg.DBName)
	fmt.Fprintln(Stdout, "\nmysql: "+q)
	_, err = conn.ExecContext(c.Ctx(), q)
	return err
}

// DropDB delete database
func (t *MySQLTool) DropDB(c *typgo.Context) error {
	conn, err := t.createAdminConn()
	if err != nil {
		return err
	}
	defer conn.Close()

	q := fmt.Sprintf("DROP DATABASE IF EXISTS `%s`", t.cfg.DBName)
	fmt.Fprintln(Stdout, "\nmysql: "+q)
	_, err = conn.ExecContext(c.Ctx(), q)
	return err
}

// MigrateDB migrate database
func (t *MySQLTool) MigrateDB(c *typgo.Context) error {
	fmt.Fprintf(Stdout, "\nmysql: Migrate '%s'\n", t.MigrationSrc)
	migration, err := t.createMigration()
	if err != nil {
		return err
	}
	defer migration.Close()
	return migration.Up()
}

// RollbackDB rollback database
func (t *MySQLTool) RollbackDB(c *typgo.Context) error {
	fmt.Fprintf(Stdout, "\nmysql: Rollback '%s'\n", t.MigrationSrc)
	migration, err := t.createMigration()
	if err != nil {
		return err
	}
	defer migration.Close()
	return migration.Down()
}

// SeedDB seed database
func (t *MySQLTool) SeedDB(c *typgo.Context) error {
	db, err := t.createConn()
	if err != nil {
		return err
	}
	defer db.Close()

	files, _ := ioutil.ReadDir(t.SeedSrc)
	for _, f := range files {
		filename := fmt.Sprintf("%s/%s", t.SeedSrc, f.Name())
		fmt.Printf("\nmysql: Seed '%s'\n", filename)
		b, _ := ioutil.ReadFile(filename)
		_, err = db.ExecContext(c.Ctx(), string(b))
		if err != nil {
			return err
		}
	}
	return nil
}

func (t *MySQLTool) createMigration() (*migrate.Migrate, error) {
	db, err := t.createConn()
	if err != nil {
		return nil, err
	}
	driver, err := mysql.WithInstance(db, &mysql.Config{})
	if err != nil {
		return nil, err
	}
	return migrate.NewWithDatabaseInstance(t.MigrationSrc, "mysql", driver)
}

func (t *MySQLTool) createConn() (*sql.DB, error) {
	return sql.Open("mysql", fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?tls=false&multiStatements=true",
		t.cfg.DBUser, t.cfg.DBPass, t.cfg.Host, t.cfg.Port, t.cfg.DBName,
	))
}

func (t *MySQLTool) createAdminConn() (*sql.DB, error) {
	return sql.Open("mysql", fmt.Sprintf(
		"root:%s@tcp(%s:%s)/?tls=false&multiStatements=true",
		t.cfg.DBPass, t.cfg.Host, t.cfg.Port,
	))
}
