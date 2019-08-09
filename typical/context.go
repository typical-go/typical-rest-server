package typical

import (
	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/typictx"
	"github.com/typical-go/typical-rest-server/app"
	"github.com/typical-go/typical-rest-server/config"
	"github.com/typical-go/typical-rest-server/typical/module/postgres"
	"github.com/typical-go/typical-rest-server/typical/module/server"

	log "github.com/sirupsen/logrus"
)

// Context instance of Context
var Context = typictx.Context{
	Name:        "Typical-RESTful-Server",
	Version:     "0.4.4",
	Description: "Example of typical and scalable RESTful API Server for Go",

	Configurations: []*typictx.Config{
		{Prefix: "APP", Spec: &config.Config{}, Description: "Application configuration"},
		{Prefix: "SERVER", Spec: &server.Config{}, Description: "Application configuration"},
		{Prefix: "PG", Spec: &postgres.Config{}, Description: "Postgres configuration"},
	},

	Application: typictx.Application{
		Action: typictx.MainAction{
			StartFunc: func(s *app.Server, cfg *config.Config) error {
				log.Info("Start the application")
				return s.Start(cfg.Address)
			},
		},
	},

	Modules: []*typictx.Module{
		server.Module(),
		postgres.Module(),
	},

	Github: &typictx.Github{
		Owner:    "typical-go",
		RepoName: "typical-rest-server",
	},
}
