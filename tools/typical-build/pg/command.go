package pg

import (
	"io"
	"log"
	"os"

	"github.com/typical-go/typical-go/pkg/execkit"
	"github.com/typical-go/typical-go/pkg/typgo"
	"github.com/typical-go/typical-rest-server/internal/generated"
	"github.com/urfave/cli/v2"

	// load migration file

	_ "github.com/golang-migrate/migrate/source/file"
)

type (
	// Command for postgres
	Command struct{}
)

var _ (typgo.Cmd) = (*Command)(nil)

// Stdout standard output
var Stdout io.Writer = os.Stdout

const (
	migrationSrc = "file://databases/pg/migration"
	seedSrc      = "databases/pg/seed"
	dockerName   = "typical-rest-server_pg01_1"
)

// Command for postgress
func (*Command) Command(sys *typgo.BuildSys) *cli.Command {
	postgresCfg, err := generated.LoadPostgresCfg()
	if err != nil {
		log.Fatal(err.Error())
	}

	u := &utility{
		PostgresCfg:  postgresCfg,
		migrationSrc: migrationSrc,
		seedSrc:      seedSrc,
	}

	return &cli.Command{
		Name:  "pg",
		Usage: "postgres utility",
		Subcommands: []*cli.Command{
			{
				Name:   "create",
				Usage:  "Create database",
				Action: sys.ExecuteFn(u.CreateDB),
			},
			{
				Name:   "drop",
				Usage:  "Drop database",
				Action: sys.ExecuteFn(u.DropDB),
			},
			{
				Name:   "migrate",
				Usage:  "Migrate database",
				Action: sys.ExecuteFn(u.MigrateDB),
			},
			{
				Name:   "rollback",
				Usage:  "Rollback database",
				Action: sys.ExecuteFn(u.RollbackDB),
			},
			{
				Name:   "seed",
				Usage:  "Seed database",
				Action: sys.ExecuteFn(u.SeedDB),
			},
			{
				Name:  "console",
				Usage: "Postgres console",
				Action: sys.ExecuteFn(
					func(c *typgo.Context) error {
						os.Setenv("PGPASSWORD", u.DBPass)
						return c.Execute(&execkit.Command{
							Name: "docker",
							Args: []string{
								"exec", "-it", dockerName,
								"psql",
								"-h", u.Host,
								"-p", u.Port,
								"-U", u.DBUser,
							},
							Stdout: os.Stdout,
							Stderr: os.Stderr,
							Stdin:  os.Stdin,
						})
					},
				),
			},
		},
	}
}
