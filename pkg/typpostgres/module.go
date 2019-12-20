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

// Config is postgres configuration
type Config struct {
	DBName   string `required:"true"`
	User     string `required:"true" default:"postgres"`
	Password string `required:"true" default:"pgpass"`
	Host     string `default:"localhost"`
	Port     int    `default:"5432"`
}

// Module of postgres
type Module struct {
	DBName string
}

// Configure the module
func (m Module) Configure() (prefix string, spec, loadFn interface{}) {
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
func (m Module) BuildCommands(c *typcore.Context) []*cli.Command {
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
func (m Module) Provide() []interface{} {
	return []interface{}{
		m.connect,
	}
}

// Prepare the module
func (m Module) Prepare() []interface{} {
	return []interface{}{
		m.ping,
	}
}

// Destroy dependencies
func (m Module) Destroy() []interface{} {
	return []interface{}{
		m.disconnect,
	}
}
