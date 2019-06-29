package ipostgres

import (
	"github.com/typical-go/typical-go/appx"
	"github.com/typical-go/typical-rest-server/typical/ext/xdb"
	"go.uber.org/dig"
	"gopkg.in/urfave/cli.v1"
)

// PostgresModule is module of postgres database
type PostgresModule struct {
	config                    interface{}
	configPrefix              string
	constructors              []interface{}
	DefaultMigrationDirectory string
	dbInfra                   appx.DBInfra
}

// NewPostgresModule return new instance of PostgresModule
func NewPostgresModule() *PostgresModule {
	return &PostgresModule{
		config:       &PGConfig{},
		configPrefix: "PG",
		constructors: []interface{}{
			LoadPostgresConfig,
			Connect,
			CreateDBInfra,
		},
	}
}

// Config to return the configuration
func (m *PostgresModule) Config() interface{} {
	return m.config
}

func (m *PostgresModule) ConfigPrefix() string {
	return m.configPrefix
}

func (m *PostgresModule) Command() cli.Command {
	return cli.Command{
		Name:      "database",
		ShortName: "db",
		Subcommands: []cli.Command{
			{Name: "create", Usage: "Create New Database", Action: m.Invoke(xdb.Create)},
			{Name: "drop", Usage: "Drop Database", Action: m.Invoke(xdb.Drop)},
			{Name: "migrate", Usage: "Migrate Database", Action: m.Invoke(xdb.Migrate)},
			{Name: "rollback", Usage: "Rollback Database", Action: m.Invoke(xdb.Rollback)},
		},
	}
}

func (m *PostgresModule) Constructors() []interface{} {
	return m.constructors
}

func (m *PostgresModule) Container() *dig.Container {
	container := dig.New()
	for _, contructor := range m.Constructors() {
		container.Provide(contructor)
	}
	return container
}

// Invoke the function with DI container
func (m *PostgresModule) Invoke(invokeFunc interface{}) interface{} {
	return func(ctx *cli.Context) error {
		container := m.Container()
		container.Provide(ctx.Args)
		return container.Invoke(invokeFunc)
	}
}
