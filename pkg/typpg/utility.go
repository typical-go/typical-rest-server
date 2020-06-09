package typpg

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"

	"github.com/golang-migrate/migrate"
	"github.com/typical-go/typical-go/pkg/execkit"
	"github.com/typical-go/typical-go/pkg/typgo"
	"github.com/urfave/cli/v2"

	_ "github.com/golang-migrate/migrate/database/postgres"
	_ "github.com/golang-migrate/migrate/source/file"
)

// Utility for PG
type Utility struct {
	Name         string
	MigrationSrc string
	SeedSrc      string
	ConfigName   string
}

// Commands of module
func (u *Utility) Commands(c *typgo.BuildCli) []*cli.Command {
	return []*cli.Command{
		{
			Name:  u.Name,
			Usage: "Postgres utility",
			Subcommands: []*cli.Command{
				u.cmdCreate(c),
				u.cmdDrop(c),
				u.cmdMigrate(c),
				u.cmdRollback(c),
				u.cmdSeed(c),
				u.cmdReset(c),
				u.cmdConsole(c),
			},
		},
	}
}

func (u *Utility) cmdCreate(c *typgo.BuildCli) *cli.Command {
	return &cli.Command{
		Name:   "create",
		Usage:  "Create New Database",
		Action: c.ActionFn(u.Name, u.create),
	}
}

func (u *Utility) create(c *typgo.Context) (err error) {
	cfg, err := u.config()
	if err != nil {
		return
	}

	conn, err := sql.Open("postgres", AdminConn(cfg))
	if err != nil {
		return
	}
	defer conn.Close()

	ctx := c.Ctx()
	if err = conn.PingContext(ctx); err != nil {
		return
	}

	c.Infof("Create database  %s", cfg.DBName)
	_, err = conn.ExecContext(ctx, fmt.Sprintf(`CREATE DATABASE "%s"`, cfg.DBName))
	return
}

func (u *Utility) cmdDrop(c *typgo.BuildCli) *cli.Command {
	return &cli.Command{
		Name:   "drop",
		Usage:  "Drop Database",
		Action: c.ActionFn(u.Name, u.drop),
	}
}

func (u *Utility) drop(c *typgo.Context) (err error) {
	cfg, err := u.config()
	if err != nil {
		return
	}

	conn, err := sql.Open("postgres", AdminConn(cfg))
	if err != nil {
		return
	}
	defer conn.Close()

	c.Infof("Drop database %s", cfg.DBName)
	_, err = conn.ExecContext(c.Ctx(), fmt.Sprintf(`DROP DATABASE IF EXISTS "%s"`, cfg.DBName))
	return
}

func (u *Utility) cmdMigrate(c *typgo.BuildCli) *cli.Command {
	return &cli.Command{
		Name:   "migrate",
		Usage:  "Migrate Database",
		Action: c.ActionFn(u.Name, u.migrate),
	}
}

func (u *Utility) migrate(c *typgo.Context) (err error) {
	cfg, err := u.config()
	if err != nil {
		return
	}

	sourceURL := "file://" + u.MigrationSrc
	c.Infof("Migrate database from source '%s'", sourceURL)
	migration, err := migrate.New(sourceURL, Conn(cfg))
	if err != nil {
		return err
	}
	defer migration.Close()
	return migration.Up()
}

func (u *Utility) cmdConsole(c *typgo.BuildCli) *cli.Command {
	return &cli.Command{
		Name:    "console",
		Aliases: []string{"c"},
		Usage:   "PostgreSQL Interactive",
		Action:  c.ActionFn(u.Name, u.console),
	}
}

func (u *Utility) console(c *typgo.Context) (err error) {
	cfg, err := u.config()
	if err != nil {
		return
	}

	os.Setenv("PGPASSWORD", cfg.Password)

	// TODO: using `docker -it` for psql

	return c.Execute(&execkit.Command{
		Name: "psql",
		Args: []string{
			"-h", cfg.Host,
			"-p", strconv.Itoa(cfg.Port),
			"-U", cfg.User,
		},
		Stdout: os.Stdout,
		Stderr: os.Stderr,
		Stdin:  os.Stdin,
	})
}

func (u *Utility) cmdRollback(c *typgo.BuildCli) *cli.Command {
	return &cli.Command{
		Name:   "rollback",
		Usage:  "Rollback Database",
		Action: c.ActionFn(u.Name, u.rollback),
	}
}

func (u *Utility) rollback(c *typgo.Context) (err error) {
	cfg, err := u.config()
	if err != nil {
		return
	}

	sourceURL := "file://" + u.MigrationSrc
	c.Infof("Migrate database from source '%s'\n", sourceURL)
	migration, err := migrate.New(sourceURL, Conn(cfg))
	if err != nil {
		return
	}
	defer migration.Close()
	return migration.Down()
}

func (u *Utility) cmdReset(c *typgo.BuildCli) *cli.Command {
	return &cli.Command{
		Name:   "reset",
		Usage:  "Reset Database",
		Action: c.ActionFn(u.Name, u.reset),
	}
}

func (u *Utility) reset(c *typgo.Context) (err error) {
	if err = u.drop(c); err != nil {
		return
	}
	if err = u.create(c); err != nil {
		return
	}
	if err = u.migrate(c); err != nil {
		return
	}
	if err = u.seed(c); err != nil {
		return
	}
	return
}

func (u *Utility) cmdSeed(c *typgo.BuildCli) *cli.Command {
	return &cli.Command{
		Name:   "seed",
		Usage:  "Data seeding",
		Action: c.ActionFn(u.Name, u.seed),
	}
}

func (u *Utility) seed(c *typgo.Context) (err error) {
	cfg, err := u.config()
	if err != nil {
		return
	}

	db, err := sql.Open("postgres", Conn(cfg))
	if err != nil {
		return
	}
	defer db.Close()

	files, _ := ioutil.ReadDir(u.SeedSrc)
	ctx := c.Ctx()
	for _, f := range files {
		sqlFile := u.SeedSrc + "/" + f.Name()
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

func (u *Utility) config() (*Config, error) {
	var cfg Config
	if err := typgo.ProcessConfig(u.ConfigName, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}
