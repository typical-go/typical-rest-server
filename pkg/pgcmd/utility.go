package pgcmd

import (
	"fmt"
	"os"

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

var _ typgo.Cmd = (*Utility)(nil)

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

// Command postgres
func (u *Utility) Command(sys *typgo.BuildSys) *cli.Command {
	return &cli.Command{
		Name:  u.Name,
		Usage: "Postgres utility",
		Subcommands: []*cli.Command{
			{
				Name:   "create",
				Usage:  "Create database",
				Action: sys.ActionFn(u.createAction("create")),
			},
			{
				Name:   "drop",
				Usage:  "Drop database",
				Action: sys.ActionFn(u.createAction("drop")),
			},
			{
				Name:   "migrate",
				Usage:  "Migrate database",
				Action: sys.ActionFn(u.createAction("migrate")),
			},
			{
				Name:   "rollback",
				Usage:  "Rollback database",
				Action: sys.ActionFn(u.createAction("rollback")),
			},
			{
				Name:   "seed",
				Usage:  "Seed database",
				Action: sys.ActionFn(u.createAction("seed")),
			},
			{
				Name:  "reset",
				Usage: "Reset database",
				Action: sys.ActionFn(typgo.Actions{
					u.createAction("drop"),
					u.createAction("create"),
					u.createAction("migrate"),
					u.createAction("seed"),
				}),
			},
			{
				Name:   "console",
				Usage:  "Postgres console",
				Action: sys.ActionFn(typgo.NewAction(u.console)),
			},
		},
	}
}

func (u *Utility) createAction(op string) typgo.Action {
	return typgo.NewAction(func(c *typgo.Context) error {
		if errMsg := u.validate(); errMsg != "" {
			return fmt.Errorf("pg-cmd: %s", errMsg)
		}

		if _, err := os.Stat(bin); os.IsNotExist(err) {
			if err := c.Execute(&execkit.GoBuild{
				Output:      bin,
				MainPackage: src,
			}); err != nil {
				return err
			}
		}
		return c.Execute(&execkit.Command{
			Name: bin,
			Args: []string{
				op,
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
