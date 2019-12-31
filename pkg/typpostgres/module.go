package typpostgres

import (
	_ "github.com/golang-migrate/migrate/database/postgres"
	_ "github.com/golang-migrate/migrate/source/file"
	_ "github.com/lib/pq"
	"github.com/typical-go/typical-go/pkg/common"
	"github.com/typical-go/typical-go/pkg/typcore"
	"github.com/urfave/cli/v2"
)

const (
	migrationSrc = "scripts/db/migration"
	seedSrc      = "scripts/db/seed"
)

// New postgres module
func New() *Module {
	return &Module{}
}

// Module of postgres
type Module struct {
	DBName string
}

// WithDBName to set database name
func (m *Module) WithDBName(dbname string) *Module {
	m.DBName = dbname
	return m
}

// Configure the module
func (m *Module) Configure() (prefix string, spec, loadFn interface{}) {
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
func (m *Module) BuildCommands(c *typcore.Context) []*cli.Command {
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
func (m *Module) Provide() []interface{} {
	return []interface{}{
		m.connect,
	}
}

// Prepare the module
func (m *Module) Prepare() []interface{} {
	return []interface{}{
		m.ping,
	}
}

// Destroy dependencies
func (m *Module) Destroy() []interface{} {
	return []interface{}{
		m.disconnect,
	}
}
