package typpostgres

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strconv"

	"github.com/golang-migrate/migrate"
	log "github.com/sirupsen/logrus"
	"github.com/typical-go/typical-go/pkg/typbuildtool"
	"github.com/urfave/cli/v2"
)

// Commands of module
func (m *Module) Commands(c *typbuildtool.Context) []*cli.Command {
	return []*cli.Command{
		{
			Name:    "postgres",
			Aliases: []string{"pg"},
			Usage:   "Postgres Database Tool",
			Subcommands: []*cli.Command{
				{
					Name:   "create",
					Usage:  "Create New Database",
					Action: m.actionFunc(c, create),
				},
				{
					Name:   "drop",
					Usage:  "Drop Database",
					Action: m.actionFunc(c, drop),
				},
				{
					Name:   "migrate",
					Usage:  "Migrate Database",
					Action: m.actionFunc(c, migrateDB),
				},
				{
					Name:   "rollback",
					Usage:  "Rollback Database",
					Action: m.actionFunc(c, rollbackDB),
				},
				{
					Name:   "seed",
					Usage:  "Data seeding",
					Action: m.actionFunc(c, seed),
				},
				{
					Name:   "reset",
					Usage:  "Reset Database",
					Action: m.actionFunc(c, reset),
				},
				{
					Name:   "console",
					Usage:  "PostgreSQL Interactive",
					Action: m.actionFunc(c, console),
				},
			},
		},
	}
}

func (m *Module) actionFunc(c *typbuildtool.Context, fn func(ctx *Context) error) cli.ActionFunc {
	return func(cliCtx *cli.Context) (err error) {
		return fn(&Context{
			Context: c,
			Cli:     cliCtx,
			Module:  m,
		})
	}
}

func reset(c *Context) (err error) {
	if err = drop(c); err != nil {
		return
	}
	if err = create(c); err != nil {
		return
	}
	if err = migrateDB(c); err != nil {
		return
	}
	if err = seed(c); err != nil {
		return
	}
	return
}

func create(c *Context) (err error) {
	var conn *sql.DB
	var cfg *Config

	if cfg, err = c.Config(); err != nil {
		return
	}

	if conn, err = sql.Open("postgres", adminDataSource(cfg)); err != nil {
		return
	}
	defer conn.Close()
	if err = conn.Ping(); err != nil {
		return
	}

	query := fmt.Sprintf(`CREATE DATABASE "%s"`, cfg.DBName)
	log.Infof("Postgres: %s", query)
	_, err = conn.Exec(query)
	return
}

func drop(c *Context) (err error) {
	var conn *sql.DB
	var cfg *Config

	if cfg, err = c.Config(); err != nil {
		return
	}

	if conn, err = sql.Open("postgres", adminDataSource(cfg)); err != nil {
		return
	}
	defer conn.Close()

	query := fmt.Sprintf(`DROP DATABASE IF EXISTS "%s"`, cfg.DBName)
	log.Infof("Postgres: %s", query)
	_, err = conn.Exec(query)
	return
}

func console(c *Context) (err error) {
	var cfg *Config
	if cfg, err = c.Config(); err != nil {
		return
	}

	os.Setenv("PGPASSWORD", cfg.Password)
	// TODO: using `docker -it` for psql
	cmd := exec.Command("psql", "-h", cfg.Host, "-p", strconv.Itoa(cfg.Port), "-U", cfg.User)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stdout
	cmd.Stdin = os.Stdin
	return cmd.Run()
}

func migrateDB(c *Context) (err error) {
	var (
		migration *migrate.Migrate
		sourceURL = "file://" + c.migrationSource
		cfg       *Config
	)

	if cfg, err = c.Config(); err != nil {
		return
	}

	log.Infof("Migrate database from source '%s'\n", sourceURL)
	if migration, err = migrate.New(sourceURL, dataSource(cfg)); err != nil {
		return err
	}
	defer migration.Close()
	return migration.Up()
}

func rollbackDB(c *Context) (err error) {
	var migration *migrate.Migrate
	var cfg *Config
	if cfg, err = c.Config(); err != nil {
		return
	}

	sourceURL := "file://" + c.migrationSource
	log.Infof("Migrate database from source '%s'\n", sourceURL)
	if migration, err = migrate.New(sourceURL, dataSource(cfg)); err != nil {
		return
	}
	defer migration.Close()
	return migration.Down()
}

func seed(c *Context) (err error) {
	var conn *sql.DB
	var cfg *Config
	if cfg, err = c.Config(); err != nil {
		return
	}
	if conn, err = sql.Open("postgres", dataSource(cfg)); err != nil {
		return
	}
	defer conn.Close()
	files, _ := ioutil.ReadDir(c.seedSource)
	for _, f := range files {
		sqlFile := c.seedSource + "/" + f.Name()
		log.Infof("Execute seed '%s'", sqlFile)
		var b []byte
		if b, err = ioutil.ReadFile(sqlFile); err != nil {
			log.Error(err.Error())
			continue
		}
		if _, err = conn.Exec(string(b)); err != nil {
			log.Error(err.Error())
			continue
		}
	}
	return
}
