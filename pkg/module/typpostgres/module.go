package typpostgres

import (
	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/docker"
	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/typictx"
	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/typienv"
)

// Module for postgres
func Module() *typictx.Module {
	return &typictx.Module{
		Name: "Postgres Database",

		Config: typictx.Config{
			Prefix: "PG",
			Spec:   &Config{},
		},

		SideEffects: []*typictx.SideEffect{
			typictx.NewSideEffect("github.com/lib/pq"),
			typictx.NewSideEffect("github.com/golang-migrate/migrate/database/postgres").ExcludeApp(),
			typictx.NewSideEffect("github.com/golang-migrate/migrate/source/file").ExcludeApp(),
		},

		OpenFunc:  openConnection,
		CloseFunc: closeConnection,
		Command: &typictx.Command{
			Name:       "postgres",
			ShortName:  "pg",
			Usage:      "Postgres Database Tool",
			BeforeFunc: typienv.LoadEnv,
			SubCommands: []*typictx.Command{
				{Name: "create", Usage: "Create New Database", ActionFunc: typictx.ActionFunction(createDB)},
				{Name: "drop", Usage: "Drop Database", ActionFunc: typictx.ActionFunction(dropDB)},
				{Name: "migrate", Usage: "Migrate Database", ActionFunc: typictx.ActionFunction(migrateDB)},
				{Name: "rollback", Usage: "Rollback Database", ActionFunc: typictx.ActionFunction(rollbackDB)},
				{Name: "console", Usage: "PostgreSQL interactive terminal", ActionFunc: typictx.ActionFunction(console)},
			},
		},

		DockerCompose: docker.NewCompose("").
			RegisterService("postgres", &docker.Service{
				Image: "postgres",
				Environment: map[string]string{
					"POSTGRES":          "${PG_USER:-postgres}",
					"POSTGRES_PASSWORD": "${PG_PASSWORD:-pgpass}",
					"PGDATA":            "/data/postgres",
				},
				Volumes: []string{
					"postgres:/data/postgres",
				},
				Ports: []string{
					"${PG_PORT:-5432}:5432",
				},
				Networks: []string{
					"postgres",
				},
				Restart: "unless-stopped",
			}).
			RegisterNetwork("postgres", &docker.Network{
				Driver: "bridge",
			}).
			RegisterVolume("postgres", nil),
	}
}
