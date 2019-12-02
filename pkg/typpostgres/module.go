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
	"github.com/typical-go/typical-go/pkg/typcfg"
	"github.com/typical-go/typical-go/pkg/typcli"
	"github.com/typical-go/typical-go/pkg/utility/envfile"
	"github.com/typical-go/typical-rest-server/pkg/typdocker"
	"github.com/urfave/cli/v2"
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
	return &module{}
}

type module struct{}

// Command of module
func (p module) Commands(c *typcli.ModuleCli) []*cli.Command {
	return []*cli.Command{
		{
			Name:    "postgres",
			Aliases: []string{"pg"},
			Usage:   "Postgres Database Tool",
			Before: func(ctx *cli.Context) error {
				return envfile.Load()
			},
			Subcommands: []*cli.Command{
				{Name: "create", Usage: "Create New Database", Action: c.Action(p.createDB)},
				{Name: "drop", Usage: "Drop Database", Action: c.Action(p.dropDB)},
				{Name: "migrate", Usage: "Migrate Database", Action: c.Action(p.migrateDB)},
				{Name: "rollback", Usage: "Rollback Database", Action: c.Action(p.rollbackDB)},
				{Name: "seed", Usage: "Database Seeding", Action: c.Action(p.seedDB)},
				{Name: "console", Usage: "PostgreSQL Interactive", Action: c.Action(p.console)},
			},
		},
	}
}

// Provide dependencies
func (p module) Provide() []interface{} {
	return []interface{}{
		p.open,
	}
}

func (p module) Prepare() []interface{} {
	return []interface{}{
		p.ping,
	}
}

func (p module) Configure() (prefix string, spec, loadFn interface{}) {
	prefix = "PG"
	spec = &Config{}
	loadFn = func(loader typcfg.Loader) (cfg Config, err error) {
		err = loader.Load(prefix, &cfg)
		return
	}
	return
}

// Destroy dependencies
func (p module) Destroy() []interface{} {
	return []interface{}{
		p.closeConnection,
	}
}

func (p module) DockerCompose() typdocker.Compose {
	return typdocker.Compose{
		Services: map[string]interface{}{
			"postgres": typdocker.Service{
				Image: "postgres",
				Environment: map[string]string{
					"POSTGRES":          "${PG_USER:-postgres}",
					"POSTGRES_PASSWORD": "${PG_PASSWORD:-pgpass}",
					"PGDATA":            "/data/postgres",
				},
				Volumes:  []string{"postgres:/data/postgres"},
				Ports:    []string{"${PG_PORT:-5432}:5432"},
				Networks: []string{"postgres"},
				Restart:  "unless-stopped",
			},
		},
		Networks: map[string]interface{}{
			"postgres": typdocker.Network{
				Driver: "bridge",
			},
		},
		Volumes: map[string]interface{}{
			"postgres": nil,
		},
	}
}

func (p module) open(cfg Config) (*sql.DB, error) {
	return sql.Open("postgres", p.dataSource(cfg))
}

func (p module) ping(db *sql.DB) error {
	log.Info("Ping to Postgres")
	return db.Ping()
}

func (module) closeConnection(db *sql.DB) error {
	return db.Close()
}

func (p module) createDB(cfg Config) (err error) {
	var conn *sql.DB
	query := fmt.Sprintf(`CREATE DATABASE "%s"`, cfg.DBName)
	fmt.Println(cfg)
	log.Infof("Postgres: %s", query)
	if conn, err = sql.Open("postgres", p.adminDataSource(cfg)); err != nil {
		return
	}
	defer conn.Close()
	fmt.Println("1")
	if err = conn.Ping(); err != nil {
		return
	}
	fmt.Println("2")
	_, err = conn.Exec(query)
	return
}

func (p module) dropDB(cfg Config) (err error) {
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

func (p module) migrateDB(cfg Config) (err error) {
	var migration *migrate.Migrate
	sourceURL := "file://" + migrationSrc
	log.Infof("Migrate database from source '%s'\n", sourceURL)
	if migration, err = migrate.New(sourceURL, p.dataSource(cfg)); err != nil {
		return err
	}
	defer migration.Close()
	return migration.Up()
}

func (p module) rollbackDB(cfg Config) (err error) {
	var migration *migrate.Migrate
	sourceURL := "file://" + migrationSrc
	log.Infof("Migrate database from source '%s'\n", sourceURL)
	if migration, err = migrate.New(sourceURL, p.dataSource(cfg)); err != nil {
		return
	}
	defer migration.Close()
	return migration.Down()
}

func (p module) seedDB(cfg Config) (err error) {
	var conn *sql.DB
	if conn, err = sql.Open("postgres", p.dataSource(cfg)); err != nil {
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

func (module) console(cfg Config) (err error) {
	os.Setenv("PGPASSWORD", cfg.Password)
	// TODO: using `docker -it` for psql
	cmd := exec.Command("psql", "-h", cfg.Host, "-p", strconv.Itoa(cfg.Port), "-U", cfg.User)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stdout
	cmd.Stdin = os.Stdin
	return cmd.Run()
}

func (module) dataSource(c Config) string {
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		c.User, c.Password, c.Host, c.Port, c.DBName)
}

func (module) adminDataSource(c Config) string {
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		c.User, c.Password, c.Host, c.Port, "template1")
}
