package postgres

import (
	"database/sql"

	log "github.com/sirupsen/logrus"
	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/typictx"
)

// Module for postgres
func Module() *typictx.Module {
	return &typictx.Module{
		Name: "Postgres Database",

		SideEffects: []*typictx.SideEffect{
			typictx.NewSideEffect("github.com/lib/pq"),
			typictx.NewSideEffect("github.com/golang-migrate/migrate/database/postgres").ExcludeApp(),
			typictx.NewSideEffect("github.com/golang-migrate/migrate/source/file").ExcludeApp(),
		},

		Command: &typictx.Command{
			Name:      "postgres",
			ShortName: "pg",
			Usage:     "Postgres Database Tool",
			SubCommands: []typictx.Command{
				{Name: "create", Usage: "Create New Database", ActionFunc: CreateDB},
				{Name: "drop", Usage: "Drop Database", ActionFunc: DropDB},
				{Name: "migrate", Usage: "Migrate Database", ActionFunc: MigrateDB},
				{Name: "rollback", Usage: "Rollback Database", ActionFunc: RollbackDB},
			},
		},

		OpenFunc: func(cfg *Config) (*sql.DB, error) {
			log.Info("Open postgres connection")
			return sql.Open(cfg.DriverName(), cfg.DataSource())
		},
		CloseFunc: func(db *sql.DB) error {
			log.Info("Close postgres connection")
			return db.Close()
		},
	}
}
