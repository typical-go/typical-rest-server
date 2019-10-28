package typpostgres

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strconv"

	"github.com/golang-migrate/migrate"
	_ "github.com/golang-migrate/migrate/database/postgres"
	_ "github.com/golang-migrate/migrate/source/file"
	"github.com/kelseyhightower/envconfig"
	_ "github.com/lib/pq"
	log "github.com/sirupsen/logrus"
	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/typictx"
	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/typienv"
	"github.com/urfave/cli"
	"go.uber.org/dig"
)

const (
	migrationSrc = "scripts/db/migration"
	seedSrc      = "scripts/db/seed"
)

// Module for postgres
func Module() interface{} {
	return &postgres{
		Name:   "Postgres",
		Config: typictx.Config{Prefix: "PG", Spec: &Config{}},
	}
}

type postgres struct {
	typictx.Config
	Name string
}

func (p postgres) CommandLine() cli.Command {
	return cli.Command{
		Name:      "postgres",
		ShortName: "pg",
		Usage:     "Postgres Database Tool",
		Before:    p.cliBefore,
		Subcommands: []cli.Command{
			{Name: "create", Usage: "Create New Database", Action: p.action(p.createDB)},
			{Name: "drop", Usage: "Drop Database", Action: p.action(p.dropDB)},
			{Name: "migrate", Usage: "Migrate Database", Action: p.action(p.migrateDB)},
			{Name: "rollback", Usage: "Rollback Database", Action: p.action(p.rollbackDB)},
			{Name: "seed", Usage: "Database Seeding", Action: p.action(p.seedDB)},
			{Name: "console", Usage: "PostgreSQL interactive terminal", Action: p.action(p.console)},
		},
	}
}

func (p postgres) Construct(c *dig.Container) (err error) {
	return c.Provide(p.openConnection)
}

func (p postgres) Destruct(c *dig.Container) (err error) {
	return c.Invoke(p.closeConnection)
}

func (p postgres) Configure() typictx.Config {
	return p.Config
}

func (p postgres) cliBefore(ctx *cli.Context) (err error) {
	return typienv.LoadEnvFile()
}

func (p postgres) action(fn interface{}) func(ctx *cli.Context) error {
	return func(ctx *cli.Context) (err error) {
		var c *dig.Container
		if c, err = p.container(); err != nil {
			return
		}
		return c.Invoke(fn)
	}
}

func (p postgres) container() (c *dig.Container, err error) {
	c = dig.New()
	err = c.Provide(p.loadConfig)
	return
}

func (p postgres) loadConfig() (cfg *Config, err error) {
	cfg = new(Config)
	err = envconfig.Process(p.Configure().Prefix, cfg)
	return
}

func (p postgres) openConnection(cfg *Config) (db *sql.DB, err error) {
	log.Info("Open postgres connection")
	if db, err = sql.Open("postgres", p.dataSource(cfg)); err != nil {
		return
	}
	err = db.Ping()
	return
}

func (postgres) closeConnection(db *sql.DB) error {
	log.Info("Close postgres connection")
	return db.Close()
}

func (p postgres) createDB(cfg *Config) (err error) {
	var conn *sql.DB
	query := fmt.Sprintf(`CREATE DATABASE "%s"`, cfg.DBName)
	log.Infof("Postgres: %s", query)
	if conn, err = sql.Open("postgres", p.adminDataSource(cfg)); err != nil {
		return
	}
	defer conn.Close()
	_, err = conn.Exec(query)
	return
}

func (p postgres) dropDB(cfg *Config) (err error) {
	var conn *sql.DB
	query := fmt.Sprintf(`DROP DATABASE IF EXISTS "%s"`, cfg.DBName)
	log.Infof("Postgres: %s", query)
	if conn, err = sql.Open("postgres", p.adminDataSource(cfg)); err != nil {
		return
	}
	defer conn.Close()
	_, err = conn.Exec(query)
	return
}

func (p postgres) migrateDB(cfg *Config) (err error) {
	var migration *migrate.Migrate
	sourceURL := "file://" + migrationSrc
	log.Infof("Migrate database from source '%s'\n", sourceURL)
	if migration, err = migrate.New(sourceURL, p.dataSource(cfg)); err != nil {
		return err
	}
	defer migration.Close()
	return migration.Up()
}

func (p postgres) rollbackDB(cfg *Config) (err error) {
	var migration *migrate.Migrate
	sourceURL := "file://" + migrationSrc
	log.Infof("Migrate database from source '%s'\n", sourceURL)
	if migration, err = migrate.New(sourceURL, p.dataSource(cfg)); err != nil {
		return
	}
	defer migration.Close()
	return migration.Down()
}

func (p postgres) seedDB(cfg *Config) (err error) {
	conn, err := sql.Open("postgres", p.dataSource(cfg))
	if err != nil {
		return
	}
	defer conn.Close()
	files, _ := ioutil.ReadDir(seedSrc)
	for _, f := range files {
		sqlFile := seedSrc + "/" + f.Name()
		log.Infof("Execute seed '%s'", sqlFile)
		var b []byte
		if b, err = ioutil.ReadFile(sqlFile); err != nil {
			return
		}
		if _, err = conn.Exec(string(b)); err != nil {
			return
		}
	}
	return
}

func (postgres) console(cfg *Config) (err error) {
	os.Setenv("PGPASSWORD", cfg.Password)
	// TODO: using `docker -it` for psql
	cmd := exec.Command("psql", "-h", cfg.Host, "-p", strconv.Itoa(cfg.Port), "-U", cfg.User)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stdout
	cmd.Stdin = os.Stdin
	return cmd.Run()
}

func (p postgres) dataSource(c *Config) string {
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		c.User, c.Password, c.Host, c.Port, c.DBName)
}

func (postgres) adminDataSource(c *Config) string {
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		c.User, c.Password, c.Host, c.Port, "template1")
}
