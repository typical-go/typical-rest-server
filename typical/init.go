package typical

import (
	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/typictx"
	"github.com/typical-go/typical-rest-server/app"
	"github.com/typical-go/typical-rest-server/app/controller"
	"github.com/typical-go/typical-rest-server/app/repository"
	"github.com/typical-go/typical-rest-server/config"
	"github.com/typical-go/typical-rest-server/typical/module"

	log "github.com/sirupsen/logrus"
)

// Context instance of Context
var Context = typictx.Context{
	Name:        "Typical-RESTful-Server",
	Version:     "0.2.2",
	Description: "Example of typical and scalable RESTful API Server for Go",

	Github: &typictx.Github{
		Owner:    "typical-go",
		RepoName: "typical-rest-server",
	},

	Configs: []typictx.Config{
		{Prefix: "APP", Spec: &config.AppConfig{}, Description: "Application configuration"},
		{Prefix: "PG", Spec: &config.PostgresConfig{}, Description: "Postgres configuration"},
	},

	AppAction: typictx.MainAction{
		StartFunc: func(s *app.Server) error {
			log.Info("Start the application")
			return s.Serve()
		},
		StopFunc: func(s *app.Server) (err error) {
			log.Info("Stop the application")
			return s.Shutdown()
		},
	},

	Constructors: []interface{}{
		app.NewServer,
		controller.NewBookController,
		controller.NewApplicationController,
		repository.NewBookRepository,
	},

	TestTargets: []string{
		"./app/controller",
		"./app/repository",
	},
	MockTargets: []string{
		"./app/repository/book_repo.go",
	},
	Modules: []*typictx.Module{
		module.NewPostgres(),
	},
}
