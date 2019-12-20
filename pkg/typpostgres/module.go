package typpostgres

import (
	_ "github.com/golang-migrate/migrate/database/postgres"
	_ "github.com/golang-migrate/migrate/source/file"
	_ "github.com/lib/pq"
	"github.com/typical-go/typical-go/pkg/typcore"
	"github.com/typical-go/typical-go/pkg/utility/common"
	"github.com/urfave/cli/v2"
)

const (
	migrationSrc = "scripts/db/migration"
	seedSrc      = "scripts/db/seed"
)

// Module of postgress
func Module(dbname string) interface{} {
	return &module{
		DBName: dbname,
	}
}

type module struct {
	DBName string
}

// Configure the module
func (m module) Configure() (prefix string, spec, loadFn interface{}) {
	prefix = "PG"
	spec = &Config{
		DBName: m.DBName,
	}
	loadFn = func(loader typcore.ConfigLoader) (cfg Config, err error) {
		err = loader.Load(prefix, &cfg)
		return
	}
	return
}

// BuildCommands of module
func (m module) BuildCommands(c *typcore.Context) []*cli.Command {
	return []*cli.Command{
		{
			Name:    "postgres",
			Aliases: []string{"pg"},
			Usage:   "Postgres Database Tool",
			Before: func(ctx *cli.Context) error {
				return common.LoadEnvFile()
			},
			Subcommands: []*cli.Command{
				m.createCmd(c),
				m.dropCmd(c),
				m.migrateCmd(c),
				m.rollbackCmd(c),
				m.seedCmd(c),
				m.resetCmd(c),
				m.consoleCmd(c),
			},
		},
	}
}

// Provide the dependencies
func (m module) Provide() []interface{} {
	return []interface{}{
		m.connect,
	}
}

// Prepare the module
func (m module) Prepare() []interface{} {
	return []interface{}{
		m.ping,
	}
}

// Destroy dependencies
func (m module) Destroy() []interface{} {
	return []interface{}{
		m.disconnect,
	}
}
