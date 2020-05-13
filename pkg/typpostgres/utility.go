package typpostgres

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strconv"

	"github.com/golang-migrate/migrate"
	"github.com/typical-go/typical-go/pkg/typgo"
	"github.com/urfave/cli/v2"

	_ "github.com/golang-migrate/migrate/database/postgres"
	_ "github.com/golang-migrate/migrate/source/file"
)

type utility struct {
	configName   string
	migrationSrc string
	seedSrc      string
}

// Utility of postgres
func Utility(s *Setting) typgo.Utility {
	if s == nil {
		s = &Setting{}
	}
	return &utility{
		configName:   GetConfigName(s),
		seedSrc:      GetSeedSrc(s),
		migrationSrc: GetMigrationSrc(s),
	}
}

// Commands of module
func (u *utility) Commands(c *typgo.BuildTool) []*cli.Command {
	return []*cli.Command{
		{
			Name:    "postgres",
			Aliases: []string{"pg"},
			Usage:   "Postgres utility",
			Subcommands: []*cli.Command{
				{
					Name:   "create",
					Usage:  "Create New Database",
					Action: c.ActionFunc("PG", u.createDB),
				},
				{
					Name:   "drop",
					Usage:  "Drop Database",
					Action: c.ActionFunc("PG", u.dropDB),
				},
				{
					Name:   "migrate",
					Usage:  "Migrate Database",
					Action: c.ActionFunc("PG", u.migrateDB),
				},
				{
					Name:   "rollback",
					Usage:  "Rollback Database",
					Action: c.ActionFunc("PG", u.rollbackDB),
				},
				{
					Name:   "seed",
					Usage:  "Data seeding",
					Action: c.ActionFunc("PG", u.seedDB),
				},
				{
					Name:   "reset",
					Usage:  "Reset Database",
					Action: c.ActionFunc("PG", u.resetDB),
				},
				{
					Name:   "console",
					Usage:  "PostgreSQL Interactive",
					Action: c.ActionFunc("PG", u.console),
				},
			},
		},
	}
}

func (u *utility) dropDB(c *typgo.Context) (err error) {
	var (
		conn *sql.DB
		cfg  *Config
	)

	if cfg, err = u.retrieveConfig(); err != nil {
		return
	}

	if conn, err = sql.Open("postgres", cfg.Admin().ConnStr()); err != nil {
		return
	}
	defer conn.Close()

	query := fmt.Sprintf(`DROP DATABASE IF EXISTS "%s"`, cfg.DBName)
	c.Infof("Postgres: %s", query)
	_, err = conn.ExecContext(c.Cli.Context, query)
	return
}

func (u *utility) createDB(c *typgo.Context) (err error) {
	var (
		conn *sql.DB
		cfg  *Config
	)

	if cfg, err = u.retrieveConfig(); err != nil {
		return
	}

	if conn, err = sql.Open("postgres", cfg.Admin().ConnStr()); err != nil {
		return
	}
	defer conn.Close()

	ctx := c.Cli.Context
	if err = conn.PingContext(ctx); err != nil {
		return
	}

	query := fmt.Sprintf(`CREATE DATABASE "%s"`, cfg.DBName)
	c.Infof("Postgres: %s", query)
	_, err = conn.ExecContext(ctx, query)
	return
}

func (u *utility) migrateDB(c *typgo.Context) (err error) {
	var (
		migration *migrate.Migrate
		cfg       *Config
	)

	if cfg, err = u.retrieveConfig(); err != nil {
		return
	}

	sourceURL := "file://" + u.migrationSrc
	c.Infof("Migrate database from source '%s'", sourceURL)
	if migration, err = migrate.New(sourceURL, cfg.ConnStr()); err != nil {
		return err
	}
	defer migration.Close()
	return migration.Up()
}

func (u *utility) console(c *typgo.Context) (err error) {
	var cfg *Config
	if cfg, err = u.retrieveConfig(); err != nil {
		return
	}

	os.Setenv("PGPASSWORD", cfg.Password)
	// TODO: using `docker -it` for psql
	cmd := exec.CommandContext(c.Cli.Context, "psql", "-h", cfg.Host, "-p", strconv.Itoa(cfg.Port), "-U", cfg.User)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stdout
	cmd.Stdin = os.Stdin
	return cmd.Run()
}

func (u *utility) rollbackDB(c *typgo.Context) (err error) {
	var (
		migration *migrate.Migrate
		cfg       *Config
	)

	if cfg, err = u.retrieveConfig(); err != nil {
		return
	}

	sourceURL := "file://" + u.migrationSrc
	c.Infof("Migrate database from source '%s'\n", sourceURL)
	if migration, err = migrate.New(sourceURL, cfg.ConnStr()); err != nil {
		return
	}
	defer migration.Close()
	return migration.Down()
}

func (u *utility) resetDB(c *typgo.Context) (err error) {

	if err = u.dropDB(c); err != nil {
		return
	}
	if err = u.createDB(c); err != nil {
		return
	}
	if err = u.migrateDB(c); err != nil {
		return
	}
	if err = u.seedDB(c); err != nil {
		return
	}
	return
}

func (u *utility) seedDB(c *typgo.Context) (err error) {
	var (
		db  *sql.DB
		cfg *Config
	)

	if cfg, err = u.retrieveConfig(); err != nil {
		return
	}

	if db, err = sql.Open("postgres", cfg.ConnStr()); err != nil {
		return
	}
	defer db.Close()

	files, _ := ioutil.ReadDir(u.seedSrc)
	ctx := c.Cli.Context
	for _, f := range files {
		sqlFile := u.seedSrc + "/" + f.Name()
		c.Infof("Execute seed '%s'", sqlFile)
		var b []byte
		if b, err = ioutil.ReadFile(sqlFile); err != nil {
			c.Warn(err.Error())
			continue
		}
		if _, err = db.ExecContext(ctx, string(b)); err != nil {
			c.Warn(err.Error())
		}
	}
	return
}

func (u *utility) retrieveConfig() (*Config, error) {
	var cfg Config
	if err := typgo.ProcessConfig(u.configName, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}
