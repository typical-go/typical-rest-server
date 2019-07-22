package module

import (
	"database/sql"

	"github.com/kelseyhightower/envconfig"
	"github.com/typical-go/typical-rest-server/config"

	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/typictx"
	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/typidb"
	log "github.com/sirupsen/logrus"
)

// NewPostgres return new instance of Module for postgrs
func NewPostgres() *typictx.Module {
	tool := typidb.NewPostgresTool()

	return &typictx.Module{
		Name:         "postgres",
		ShortName:    "pg",
		Usage:        "Postgres Database Module",
		ConfigPrefix: "PG",
		Config:       &config.PostgresConfig{},
		SideEffects: []*typictx.SideEffect{
			typictx.NewSideEffect("github.com/lib/pq"),
			typictx.NewSideEffect("github.com/golang-migrate/migrate/database/postgres").ExcludeApp(),
			typictx.NewSideEffect("github.com/golang-migrate/migrate/source/file").ExcludeApp(),
		},
		Commands: []typictx.Command{
			{Name: "create", Usage: "Create New Database", ActionFunc: tool.CreateDB},
			{Name: "drop", Usage: "Drop Database", ActionFunc: tool.DropDB},
			{Name: "migrate", Usage: "Migrate Database", ActionFunc: tool.MigrateDB},
			{Name: "rollback", Usage: "Rollback Database", ActionFunc: tool.RollbackDB},
		},
		Constructors: []interface{}{
			func() (typidb.Config, error) {
				var cfg config.PostgresConfig
				err := envconfig.Process("PG", &cfg)
				return &cfg, err
			},
		},
		OpenFunc: func(cfg typidb.Config) (*sql.DB, error) {
			log.Info("Open postgres connection")
			return sql.Open(cfg.DriverName(), cfg.DataSource())
		},
		CloseFunc: func(db *sql.DB) error {
			log.Info("Close postgres connection")
			return db.Close()
		},
		StatusFunc: func(db *sql.DB) error {
			return db.Ping()
		},
	}
}
