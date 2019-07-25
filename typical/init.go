package typical

import (
	log "github.com/sirupsen/logrus"
	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/typictx"
	"github.com/typical-go/typical-rest-server/app"
	"github.com/typical-go/typical-rest-server/app/controller"
	"github.com/typical-go/typical-rest-server/app/repository"
	"github.com/typical-go/typical-rest-server/config"
	"github.com/typical-go/typical-rest-server/typical/module"
)

// Context instance of Context
var Context typictx.Context

func init() {
	Context = typictx.Context{
		Name:        "Typical-RESTful-Server",
		Version:     "0.2.1",
		Description: "Example of typical and scalable RESTful API Server for Go",

		Github: &typictx.Github{
			Owner:    "typical-go",
			RepoName: "typical-rest-server",
		},

		Constructors: []interface{}{
			config.LoadConfig,
			config.GetAppConfig,
			config.GetPgConfig,
			config.GetDBToolConfig,
		},

		AppModule: typictx.TypiApp{
			Config:       &config.AppConfig{},
			ConfigPrefix: "APP",
			Constructors: []interface{}{
				app.NewServer,
				controller.NewBookController,
				controller.NewApplicationController,
				repository.NewBookRepository,
			},
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
			TestTargets: []string{
				"./app/controller",
				"./app/repository",
			},
			MockTargets: []string{
				"./app/repository/book_repo.go",
			},
		},
		Modules: []*typictx.Module{
			module.NewPostgres(),
		},
	}

}
