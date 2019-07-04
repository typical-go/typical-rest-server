package xpostgres

import (
	"github.com/kelseyhightower/envconfig"
	"github.com/typical-go/typical-go/appx"
	"github.com/typical-go/typical-rest-server/typical/appctx"
	"github.com/typical-go/typical-rest-server/typical/ext/xdb"
	"gopkg.in/urfave/cli.v1"
)

// PostgresModule is module of postgres database
type PostgresModule struct {
	appctx.DependencyInjection
	appctx.ModuleDetail
	DefaultMigrationDirectory string
	dbInfra                   appx.DBInfra
}

// NewModule return new instance of PostgresModule
func NewModule() *PostgresModule {
	return &PostgresModule{
		ModuleDetail: *appctx.NewModuleDetail().
			SetName("Postgres").
			SetShortName("pg").
			SetConfigPrefix("PG").
			SetConfig(&Config{}),

		DependencyInjection: appctx.NewDependencyInjection(
			Connect,
			NewDBInfra,
		),
	}
}

// LoadFunc to load the configuration
func (m *PostgresModule) LoadFunc() interface{} {
	return func() (config Config, err error) {
		err = envconfig.Process(m.ConfigPrefix(), &config)
		return
	}
}

// Command of the module
func (m *PostgresModule) Command() cli.Command {
	return cli.Command{
		Name:      m.Name(),
		ShortName: m.ShortName(),
		Usage:     "Postgres database tool",
		Subcommands: []cli.Command{
			{Name: "create", Usage: "Create New Database", Action: m.InvokeFunction(xdb.Create)},
			{Name: "drop", Usage: "Drop Database", Action: m.InvokeFunction(xdb.Drop)},
			{Name: "migrate", Usage: "Migrate Database", Action: m.InvokeFunction(xdb.Migrate)},
			{Name: "rollback", Usage: "Rollback Database", Action: m.InvokeFunction(xdb.Rollback)},
		},
	}
}
