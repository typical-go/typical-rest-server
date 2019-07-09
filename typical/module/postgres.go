package module

import (
	"database/sql"

	"github.com/kelseyhightower/envconfig"
	"github.com/typical-go/typical-rest-server/config"
	"github.com/typical-go/typical-rest-server/typical/appctx"
	"github.com/typical-go/typical-rest-server/typical/typidb"
	"gopkg.in/urfave/cli.v1"
)

// NewPostgres return new instance of Module for postgrs
func NewPostgres() *appctx.Module {
	tool := typidb.NewPostgresTool()

	m := &appctx.Module{
		Name:         "postgres",
		ShortName:    "pg",
		Usage:        "Postgres Database Module",
		ConfigPrefix: "PG",
		Config:       &config.PostgresConfig{},
	}
	m.Command = &cli.Command{
		Name:      m.Name,
		ShortName: m.ShortName,
		Usage:     m.Usage,
		Subcommands: []cli.Command{
			{Name: "create", Usage: "Create New Database", Action: m.Invoke(tool.CreateDB)},
			{Name: "drop", Usage: "Drop Database", Action: m.Invoke(tool.DropDB)},
			{Name: "migrate", Usage: "Migrate Database", Action: m.Invoke(tool.MigrateDB)},
			{Name: "rollback", Usage: "Rollback Database", Action: m.Invoke(tool.RollbackDB)},
		},
	}
	m.LoadConfigFunc = func() (typidb.Config, error) {
		var cfg config.PostgresConfig
		err := envconfig.Process(m.ConfigPrefix, &cfg)
		return &cfg, err
	}
	m.OpenFunc = func(cfg typidb.Config) (*sql.DB, error) {
		return sql.Open(cfg.DriverName(), cfg.DataSource())
	}
	return m
}
