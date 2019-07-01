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
	config                    interface{}
	configPrefix              string
	DefaultMigrationDirectory string
	dbInfra                   appx.DBInfra
}

// NewModule return new instance of PostgresModule
func NewModule() *PostgresModule {
	return &PostgresModule{
		DependencyInjection: appctx.NewDependencyInjection(
			LoadPostgresConfig,
			Connect,
			CreateDBInfra,
		),
		config:       &PGConfig{},
		configPrefix: "PG",
	}
}

// Config to return the configuration
func (m *PostgresModule) Config() interface{} {
	return m.config
}

// ConfigPrefix of the module
func (m *PostgresModule) ConfigPrefix() string {
	return m.configPrefix
}

// Command of the module
func (m *PostgresModule) Command() cli.Command {
	return cli.Command{
		Name:      "postgres",
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
