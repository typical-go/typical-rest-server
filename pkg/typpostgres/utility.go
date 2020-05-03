package typpostgres

import (
	"github.com/typical-go/typical-go/pkg/typbuildtool"
	"github.com/typical-go/typical-go/pkg/typcfg"
	"github.com/urfave/cli/v2"
)

var (
	// DefaultMigrationSource is default migration source for postgres
	DefaultMigrationSource = "scripts/db/migration"

	// DefaultSeedSource is default seed source for postgres
	DefaultSeedSource = "scripts/db/seed"
)

// Utility return new instance of PostgresUtility
func Utility() typbuildtool.Utility {
	return typbuildtool.NewUtility(Commands).
		Configure(&typcfg.Configuration{
			Name: DefaultConfigName,
			Spec: DefaultConfig,
		})
}

// Commands of module
func Commands(c *typbuildtool.Context) []*cli.Command {
	return []*cli.Command{
		{
			Name:    "postgres",
			Aliases: []string{"pg"},
			Usage:   "Postgres utility",
			Subcommands: []*cli.Command{
				cmdCreateDB(c),
				cmdDropDB(c),
				cmdMigrateDB(c),
				cmdRollbackDB(c),
				cmdSeedDB(c),
				cmdResetDB(c),
				cmdConsole(c),
			},
		},
	}
}
