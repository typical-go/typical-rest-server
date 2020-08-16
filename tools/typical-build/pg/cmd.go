package pg

import (
	"log"
	"os"

	"github.com/typical-go/typical-go/pkg/execkit"
	"github.com/typical-go/typical-go/pkg/typgo"
	"github.com/typical-go/typical-rest-server/internal/generated"
	"github.com/urfave/cli/v2"
)

type (
	// Cmd command for postgres
	Cmd struct {
	}
)

var _ (typgo.Cmd) = (*Cmd)(nil)

// Command for postgress
func (*Cmd) Command(sys *typgo.BuildSys) *cli.Command {
	postgresCfg, err := generated.LoadPostgresCfg()
	if err != nil {
		log.Fatal(err.Error())
	}

	tool := &Tool{
		PostgresCfg:  postgresCfg,
		migrationSrc: "file://databases/pg/migration",
		seedSrc:      "databases/pg/seed",
	}

	return &cli.Command{
		Name:    "postgres",
		Aliases: []string{"pg"},
		Usage:   "postgres utility",
		Subcommands: []*cli.Command{
			{
				Name:   "create",
				Usage:  "Create database",
				Action: sys.ActionFn(typgo.NewAction(tool.CreateDB)),
			},
			{
				Name:   "drop",
				Usage:  "Drop database",
				Action: sys.ActionFn(typgo.NewAction(tool.DropDB)),
			},
			{
				Name:   "migrate",
				Usage:  "Migrate database",
				Action: sys.ActionFn(typgo.NewAction(tool.MigrateDB)),
			},
			{
				Name:   "rollback",
				Usage:  "Rollback database",
				Action: sys.ActionFn(typgo.NewAction(tool.RollbackDB)),
			},
			{
				Name:   "seed",
				Usage:  "Seed database",
				Action: sys.ActionFn(typgo.NewAction(tool.SeedDB)),
			},
			{
				Name:  "reset",
				Usage: "Reset database",
				Action: sys.ActionFn(typgo.NewAction(
					func(c *typgo.Context) error {
						if err := tool.DropDB(c); err != nil {
							return err
						}
						tool.CreateDB(c)
						tool.MigrateDB(c)
						tool.SeedDB(c)
						return nil
					})),
			},
			{
				Name:  "console",
				Usage: "Postgres console",
				Action: sys.ActionFn(typgo.NewAction(
					func(c *typgo.Context) error {
						os.Setenv("PGPASSWORD", tool.Pass)

						return c.Execute(&execkit.Command{
							Name: "docker",
							Args: []string{
								"exec", "-it", "typical-rest-server_pg01_1",
								"psql",
								"-h", tool.Host,
								"-p", tool.Port,
								"-U", tool.User,
							},
							Stdout: os.Stdout,
							Stderr: os.Stderr,
							Stdin:  os.Stdin,
						})
					})),
			},
		},
	}
}
