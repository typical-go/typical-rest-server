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
	_ "github.com/lib/pq"
	log "github.com/sirupsen/logrus"
	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/typicli"
	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/typiobj"
	"github.com/urfave/cli"
)

const (
	migrationSrc = "scripts/db/migration"
	seedSrc      = "scripts/db/seed"
)

// Config is postgres configuration
type Config struct {
	DBName   string `required:"true" default:"typical-rest"`
	User     string `required:"true" default:"postgres"`
	Password string `required:"true" default:"pgpass"`
	Host     string `default:"localhost"`
	Port     int    `default:"5432"`
}

// Module for postgres
func Module() interface{} {
	return &postgresModule{
		Name: "Postgres",
		Configuration: typiobj.Configuration{
			Prefix: "PG",
			Spec:   &Config{},
		},
	}
}

type postgresModule struct {
	typiobj.Configuration
	Name string
}

// Command of module
func (p postgresModule) Command() cli.Command {
	return cli.Command{
		Name:      "postgres",
		ShortName: "pg",
		Usage:     "Postgres Database Tool",
		Before:    typicli.LoadEnvFile,
		Subcommands: []cli.Command{
			{Name: "create", Usage: "Create New Database", Action: typicli.Action(p, p.createDB)},
			{Name: "drop", Usage: "Drop Database", Action: typicli.Action(p, p.dropDB)},
			{Name: "migrate", Usage: "Migrate Database", Action: typicli.Action(p, p.migrateDB)},
			{Name: "rollback", Usage: "Rollback Database", Action: typicli.Action(p, p.rollbackDB)},
			{Name: "seed", Usage: "Database Seeding", Action: typicli.Action(p, p.seedDB)},
			{Name: "console", Usage: "PostgreSQL Interactive", Action: typicli.Action(p, p.console)},
		},
	}
}

// Provide dependencies
func (p postgresModule) Provide() []interface{} {
	return []interface{}{
		p.loadConfig,
		p.openConnection,
	}
}

// Destroy dependencies
func (p postgresModule) Destroy() []interface{} {
	return []interface{}{
		p.closeConnection,
	}
}

func (p postgresModule) loadConfig() (cfg *Config, err error) {
	err = p.Configuration.Load()
	cfg = p.Configuration.Spec.(*Config)
	return
}

func (p postgresModule) openConnection(cfg *Config) (db *sql.DB, err error) {
	log.Info("Open postgres connection")
	if db, err = sql.Open("postgres", p.dataSource(cfg)); err != nil {
		return
	}
	err = db.Ping()
	return
}

func (postgresModule) closeConnection(db *sql.DB) error {
	fmt.Println("Close postgres connection")
	return db.Close()
}

func (p postgresModule) createDB(cfg *Config) (err error) {
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

func (p postgresModule) dropDB(cfg *Config) (err error) {
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

func (p postgresModule) migrateDB(cfg *Config) (err error) {
	var migration *migrate.Migrate
	sourceURL := "file://" + migrationSrc
	log.Infof("Migrate database from source '%s'\n", sourceURL)
	if migration, err = migrate.New(sourceURL, p.dataSource(cfg)); err != nil {
		return err
	}
	defer migration.Close()
	return migration.Up()
}

func (p postgresModule) rollbackDB(cfg *Config) (err error) {
	var migration *migrate.Migrate
	sourceURL := "file://" + migrationSrc
	log.Infof("Migrate database from source '%s'\n", sourceURL)
	if migration, err = migrate.New(sourceURL, p.dataSource(cfg)); err != nil {
		return
	}
	defer migration.Close()
	return migration.Down()
}

func (p postgresModule) seedDB(cfg *Config) (err error) {
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

func (postgresModule) console(cfg *Config) (err error) {
	os.Setenv("PGPASSWORD", cfg.Password)
	// TODO: using `docker -it` for psql
	cmd := exec.Command("psql", "-h", cfg.Host, "-p", strconv.Itoa(cfg.Port), "-U", cfg.User)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stdout
	cmd.Stdin = os.Stdin
	return cmd.Run()
}

func (postgresModule) dataSource(c *Config) string {
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		c.User, c.Password, c.Host, c.Port, c.DBName)
}

func (postgresModule) adminDataSource(c *Config) string {
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		c.User, c.Password, c.Host, c.Port, "template1")
}
