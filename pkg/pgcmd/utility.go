package pgcmd

import (
	"fmt"
	"os"

	"github.com/typical-go/typical-go/pkg/buildkit"
	"github.com/typical-go/typical-go/pkg/execkit"
	"github.com/typical-go/typical-go/pkg/typgo"
	"github.com/urfave/cli/v2"
)

const (
	bin = ".typical-tmp/bin/pg-tool"
	src = "github.com/typical-go/typical-rest-server/tools/pg-tool"
)

// Utility for postgres
type Utility struct {
	Name         string
	HostEnv      string
	PortEnv      string
	UserEnv      string
	PasswordEnv  string
	DBNameEnv    string
	MigrationSrc string
	SeedSrc      string
}

var _ typgo.Utility = (*Utility)(nil)

func (u *Utility) validate() string {
	if u.Name == "" {
		return "missing name"
	}
	if u.HostEnv == "" {
		return "missing HostEnv"
	}
	if u.PortEnv == "" {
		return "missing PortEnv"
	}
	if u.PasswordEnv == "" {
		return "missing PasswordEnv"
	}
	if u.DBNameEnv == "" {
		return "missing DBNameEnv"
	}
	if u.MigrationSrc == "" {
		return "missing MigrationSrc"
	}
	if u.SeedSrc == "" {
		return "missing SeedSrc"
	}
	return ""
}

// Commands list
func (u *Utility) Commands(c *typgo.BuildCli) ([]*cli.Command, error) {
	if errMsg := u.validate(); errMsg != "" {
		return nil, fmt.Errorf("pg-cmd: %s", errMsg)
	}

	return []*cli.Command{
		{
			Name:  u.Name,
			Usage: "Postgres utility",
			Subcommands: []*cli.Command{
				{
					Name:  "create",
					Usage: "Create database",
					Action: c.ActionFn("PG", func(c *typgo.Context) error {
						return u.execute(c, "create")
					}),
				},
				{
					Name:  "drop",
					Usage: "Drop database",
					Action: c.ActionFn("PG", func(c *typgo.Context) error {
						return u.execute(c, "drop")
					}),
				},
				{
					Name:  "migrate",
					Usage: "Migrate database",
					Action: c.ActionFn("PG", func(c *typgo.Context) error {
						return u.execute(c, "migrate")
					}),
				},
				{
					Name:  "rollback",
					Usage: "Rollback database",
					Action: c.ActionFn("PG", func(c *typgo.Context) error {
						return u.execute(c, "rollback")
					}),
				},
				{
					Name:  "seed",
					Usage: "Seed database",
					Action: c.ActionFn("PG", func(c *typgo.Context) error {
						return u.execute(c, "rollback")
					}),
				},
				{
					Name:  "reset",
					Usage: "Reset database",
					Action: c.ActionFn("PG", func(c *typgo.Context) error {
						if err := u.execute(c, "drop"); err != nil {
							return err
						}
						if err := u.execute(c, "create"); err != nil {
							return err
						}
						if err := u.execute(c, "migrate"); err != nil {
							return err
						}
						if err := u.execute(c, "seed"); err != nil {
							return err
						}
						return nil
					}),
				},
				{
					Name:  "console",
					Usage: "Postgres console",
					Action: c.ActionFn("PG", func(c *typgo.Context) error {
						return u.console(c)
					}),
				},
			},
		},
	}, nil
}

func (u *Utility) execute(c *typgo.Context, action string) error {
	if _, err := os.Stat(bin); os.IsNotExist(err) {
		if err := c.Execute(&buildkit.GoBuild{
			Out:    bin,
			Source: src,
		}); err != nil {
			return err
		}
	}
	return u.pgTool(c, action)
}

func (u *Utility) pgTool(c *typgo.Context, action string) error {
	return c.Execute(&execkit.Command{
		Name: bin,
		Args: []string{
			action,
			"-host=" + os.Getenv(u.HostEnv),
			"-port=" + os.Getenv(u.PortEnv),
			"-user=" + os.Getenv(u.UserEnv),
			"-password=" + os.Getenv(u.PasswordEnv),
			"-db-name=" + os.Getenv(u.DBNameEnv),
			"-migration-src=" + u.MigrationSrc,
			"-seed-src=" + u.SeedSrc,
		},
		Stdout: os.Stdout,
		Stderr: os.Stderr,
		Stdin:  os.Stdin,
	})
}

// Console postgrs
func (u *Utility) console(c *typgo.Context) error {
	os.Setenv("PGPASSWORD", os.Getenv(u.PasswordEnv))

	// TODO: using `docker -it` for psql

	return c.Execute(&execkit.Command{
		Name: "psql",
		Args: []string{
			"-h", os.Getenv(u.HostEnv),
			"-p", os.Getenv(u.PortEnv),
			"-U", os.Getenv(u.UserEnv),
		},
		Stdout: os.Stdout,
		Stderr: os.Stderr,
		Stdin:  os.Stdin,
	})
}
