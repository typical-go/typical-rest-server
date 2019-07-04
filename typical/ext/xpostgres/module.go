package xpostgres

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/golang-migrate/migrate"
	"github.com/kelseyhightower/envconfig"
	"github.com/typical-go/typical-rest-server/typical/appctx"
	"gopkg.in/urfave/cli.v1"
)

// PostgresModule is module of postgres database
type PostgresModule struct {
	appctx.DependencyInjection
	appctx.ModuleDetail
	DefaultMigrationDirectory string
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
			{Name: "create", Usage: "Create New Database", Action: m.InvokeFunction(createDB)},
			{Name: "drop", Usage: "Drop Database", Action: m.InvokeFunction(dropDB)},
			{Name: "migrate", Usage: "Migrate Database", Action: m.InvokeFunction(migrateDB)},
			{Name: "rollback", Usage: "Rollback Database", Action: m.InvokeFunction(rollbackDB)},
		},
	}
}

// Create database
func createDB(config Config) (err error) {
	conn, err := sql.Open("postgres", config.ConnectionStringNoDB())
	if err != nil {
		return
	}
	defer conn.Close()

	query := fmt.Sprintf(`CREATE DATABASE "%s"`, config.DbName)
	fmt.Println(query)
	_, err = conn.Exec(query)
	return
}

// Drop database
func dropDB(config Config) (err error) {
	conn, err := sql.Open("postgres", config.ConnectionStringNoDB())
	if err != nil {
		return
	}
	defer conn.Close()

	query := fmt.Sprintf(`DROP DATABASE IF EXISTS "%s"`, config.DbName)
	fmt.Println(query)
	_, err = conn.Exec(query)
	return
}

// Migrate database
func migrateDB(config Config) error {
	source := config.MigrationSource
	log.Printf("Migrate database from source '%s'\n", source)

	migration, err := migrate.New(source, config.ConnectionString())
	if err != nil {
		return err
	}
	defer migration.Close()
	return migration.Up()
}

// Rollback database
func rollbackDB(config Config) error {
	source := config.MigrationSource
	log.Printf("Migrate database from source '%s'\n", source)

	migration, err := migrate.New(source, config.ConnectionString())
	if err != nil {
		return err
	}
	defer migration.Close()
	return migration.Down()
}
