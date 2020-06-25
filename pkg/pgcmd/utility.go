package pgcmd

import (
	"os"

	"github.com/typical-go/typical-go/pkg/buildkit"
	"github.com/typical-go/typical-go/pkg/typgo"
	"github.com/urfave/cli/v2"
)

// Utility for postgres
type Utility struct {
	HostEnv      string
	PortEnv      string
	UserEnv      string
	PasswordEnv  string
	DBNameEnv    string
	MigrationSrc string
	SeedSrc      string
}

var _ typgo.Utility = (*Utility)(nil)

// Commands list
func (u *Utility) Commands(c *typgo.BuildCli) ([]*cli.Command, error) {
	return []*cli.Command{
		{
			Name:  "pg",
			Usage: "postgres utility",
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
						return c.Execute(Console(&Param{
							Host:     os.Getenv(u.HostEnv),
							Port:     os.Getenv(u.PortEnv),
							User:     os.Getenv(u.UserEnv),
							Password: os.Getenv(u.PasswordEnv),
						}))
					}),
				},
			},
		},
	}, nil
}

func (u *Utility) execute(c *typgo.Context, action string) error {
	bin := ".typical-tmp/bin/pgutil"
	if _, err := os.Stat(bin); os.IsNotExist(err) {
		if err := c.Execute(&buildkit.GoBuild{
			Out:    bin,
			Source: "github.com/typical-go/typical-rest-server/tools/pg-tool",
		}); err != nil {
			return err
		}
	}

	return c.Execute(PgTool(&Param{
		Name:         bin,
		Action:       action,
		Host:         os.Getenv(u.HostEnv),
		Port:         os.Getenv(u.PortEnv),
		User:         os.Getenv(u.UserEnv),
		Password:     os.Getenv(u.PasswordEnv),
		DBName:       os.Getenv(u.DBNameEnv),
		MigrationSrc: u.MigrationSrc,
		SeedSrc:      u.SeedSrc,
	}))
}
