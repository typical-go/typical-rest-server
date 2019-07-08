package mpostgres

import (
	"database/sql"

	"github.com/kelseyhightower/envconfig"
	"github.com/typical-go/typical-rest-server/typical/appctx"
	"gopkg.in/urfave/cli.v1"
)

// New return new instance of Module for postgrs
func New() *appctx.Module {
	m := &appctx.Module{
		Name:         "postgres",
		ShortName:    "pg",
		Usage:        "Postgres Database Module",
		ConfigPrefix: "PG",
		Config:       &Config{},
	}
	m.Command = &cli.Command{
		Name:      m.Name,
		ShortName: m.ShortName,
		Usage:     m.Usage,
		Subcommands: []cli.Command{
			{Name: "create", Usage: "Create New Database", Action: m.Invoke(createDB)},
			{Name: "drop", Usage: "Drop Database", Action: m.Invoke(dropDB)},
			{Name: "migrate", Usage: "Migrate Database", Action: m.Invoke(migrateDB)},
			{Name: "rollback", Usage: "Rollback Database", Action: m.Invoke(rollbackDB)},
		},
	}
	m.LoadConfigFunc = func() (config Config, err error) {
		err = envconfig.Process(m.ConfigPrefix, &config)
		return
	}
	m.OpenFunc = func(conf Config) (*sql.DB, error) {
		return sql.Open("postgres", conf.ConnectionString())
	}
	return m
}
