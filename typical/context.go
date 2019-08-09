package typical

import (
	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/typictx"
	"github.com/typical-go/typical-rest-server/app"
	"github.com/typical-go/typical-rest-server/config"
	"github.com/typical-go/typical-rest-server/typical/module/postgres"

	log "github.com/sirupsen/logrus"
)

// Context instance of Context
var Context = typictx.Context{
	Name:        "Typical-RESTful-Server",
	Version:     "0.4.4",
	Description: "Example of typical and scalable RESTful API Server for Go",

	Configurations: []*typictx.Config{
		{Prefix: "APP", Spec: &config.AppConfig{}, Description: "Application configuration"},
		{Prefix: "PG", Spec: &postgres.Config{}, Description: "Postgres configuration"},
	},

	Application: typictx.Application{
		Action: typictx.MainAction{
			StartFunc: func(s *app.Server) error {
				log.Info("Start the application")
				return s.Serve()
			},
			StopFunc: func(s *app.Server) (err error) {
				log.Info("Stop the application")
				return s.Shutdown()
			},
		},
	},

	Modules: []*typictx.Module{
		postgres.Module(),
	},

	Github: &typictx.Github{
		Owner:    "typical-go",
		RepoName: "typical-rest-server",
	},
}
