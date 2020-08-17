package pg

import (
	"database/sql"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"

	"github.com/typical-go/typical-go/pkg/execkit"
	"github.com/typical-go/typical-go/pkg/typgo"
	"github.com/typical-go/typical-rest-server/internal/generated"
	"github.com/typical-go/typical-rest-server/internal/infra"
	"github.com/urfave/cli/v2"

	// load migration file
	"github.com/golang-migrate/migrate"
	"github.com/golang-migrate/migrate/database/postgres"
	_ "github.com/golang-migrate/migrate/source/file"
)

type (
	// Utility for postgres
	Utility struct{}
	utility struct {
		*infra.PostgresCfg
		migrationSrc string
		seedSrc      string
	}
)

var _ (typgo.Cmd) = (*Utility)(nil)

// Stdout standard output
var Stdout io.Writer = os.Stdout

const (
	migrationSrc = "file://databases/pg/migration"
	seedSrc      = "databases/pg/seed"
	dockerName   = "typical-rest-server_pg01_1"
)

// Command for postgress
func (*Utility) Command(sys *typgo.BuildSys) *cli.Command {
	postgresCfg, err := generated.LoadPostgresCfg()
	if err != nil {
		log.Fatal(err.Error())
	}

	u := &utility{
		PostgresCfg:  postgresCfg,
		migrationSrc: migrationSrc,
		seedSrc:      seedSrc,
	}

	return &cli.Command{
		Name:  "pg",
		Usage: "postgres utility",
		Subcommands: []*cli.Command{
			{
				Name:   "create",
				Usage:  "Create database",
				Action: sys.ExecuteFn(u.CreateDB),
			},
			{
				Name:   "drop",
				Usage:  "Drop database",
				Action: sys.ExecuteFn(u.DropDB),
			},
			{
				Name:   "migrate",
				Usage:  "Migrate database",
				Action: sys.ExecuteFn(u.MigrateDB),
			},
			{
				Name:   "rollback",
				Usage:  "Rollback database",
				Action: sys.ExecuteFn(u.RollbackDB),
			},
			{
				Name:   "seed",
				Usage:  "Seed database",
				Action: sys.ExecuteFn(u.SeedDB),
			},
			{
				Name:  "console",
				Usage: "Postgres console",
				Action: sys.ExecuteFn(
					func(c *typgo.Context) error {
						os.Setenv("PGPASSWORD", u.DBPass)
						return c.Execute(&execkit.Command{
							Name: "docker",
							Args: []string{
								"exec", "-it", dockerName,
								"psql",
								"-h", u.Host,
								"-p", u.Port,
								"-U", u.DBUser,
							},
							Stdout: os.Stdout,
							Stderr: os.Stderr,
							Stdin:  os.Stdin,
						})
					},
				),
			},
		},
	}
}

//
// utility
//

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

	fmt.Printf("\npg: Seed '%s'\n", t.seedSrc)
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
