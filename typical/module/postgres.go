package module

import (
	"database/sql"
	"log"

	"github.com/kelseyhightower/envconfig"
	"github.com/typical-go/typical-rest-server/config"

	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/typictx"
	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/typidb"
)

// NewPostgres return new instance of Module for postgrs
func NewPostgres() *typictx.Module {
	tool := typidb.NewPostgresTool()

	m := &typictx.Module{
		Name:         "postgres",
		ShortName:    "pg",
		Usage:        "Postgres Database Module",
		ConfigPrefix: "PG",
		Config:       &config.PostgresConfig{},
		SideEffects: []string{
			"github.com/lib/pq",
		},
		TaskToolSideEffects: []string{
			"github.com/golang-migrate/migrate/database/postgres",
			"github.com/golang-migrate/migrate/source/file",
		},
	}
	m.Commands = []typictx.Command{
		{Name: "create", Usage: "Create New Database", ActionFunc: tool.CreateDB},
		{Name: "drop", Usage: "Drop Database", ActionFunc: tool.DropDB},
		{Name: "migrate", Usage: "Migrate Database", ActionFunc: tool.MigrateDB},
		{Name: "rollback", Usage: "Rollback Database", ActionFunc: tool.RollbackDB},
	}
	m.Constructors = []interface{}{
		func() (typidb.Config, error) {
			var cfg config.PostgresConfig
			err := envconfig.Process(m.ConfigPrefix, &cfg)
			return &cfg, err
		},
	}
	m.OpenFunc = func(cfg typidb.Config) (*sql.DB, error) {
		log.Println("Open postgres connection")
		return sql.Open(cfg.DriverName(), cfg.DataSource())
	}
	m.CloseFunc = func(db *sql.DB) error {
		log.Println("Close postgres connection")
		return db.Close()
	}
	return m
}
