package xpostgres

import (
	"github.com/typical-go/typical-go/appx"
	"github.com/typical-go/typical-rest-server/typical/appctx"
	"github.com/typical-go/typical-rest-server/typical/ext/xdb"
	"gopkg.in/urfave/cli.v1"
)

// PostgresModule is module of postgres database
type PostgresModule struct {
	appctx.DependencyInjection
	appctx.ConfigLoader
	name                      string
	DefaultMigrationDirectory string
	dbInfra                   appx.DBInfra
}

// NewModule return new instance of PostgresModule
func NewModule() *PostgresModule {
	return &PostgresModule{
		DependencyInjection: appctx.NewDependencyInjection(
			Connect,
			CreateDBInfra,
		),
		ConfigLoader: ConfigLoader{
			ConfigDetail: appctx.NewConfigDetail("PG", &Config{}),
		},
		name: "Postgres",
	}
}

// Name of the module
func (m *PostgresModule) Name() string {
	return m.name
}

// Command of the module
func (m *PostgresModule) Command() cli.Command {
	return cli.Command{
		Name:      m.Name(),
		ShortName: "pg",
		Usage:     "Postgres database tool",
		Subcommands: []cli.Command{
			{Name: "create", Usage: "Create New Database", Action: m.InvokeFunction(xdb.Create)},
			{Name: "drop", Usage: "Drop Database", Action: m.InvokeFunction(xdb.Drop)},
			{Name: "migrate", Usage: "Migrate Database", Action: m.InvokeFunction(xdb.Migrate)},
			{Name: "rollback", Usage: "Rollback Database", Action: m.InvokeFunction(xdb.Rollback)},
		},
	}
}
