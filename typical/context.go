package typical

import (
	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/docker"
	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/typictx"
	"github.com/typical-go/typical-rest-server/app"
	"github.com/typical-go/typical-rest-server/app/config"
	"github.com/typical-go/typical-rest-server/typical/module/typpostgres"
	"github.com/typical-go/typical-rest-server/typical/module/typserver"

	log "github.com/sirupsen/logrus"
)

// Context instance of Context
var Context = typictx.Context{
	Name:        "Typical-RESTful-Server",
	Version:     "0.6.0",
	Description: "Example of typical and scalable RESTful API Server for Go",

	Modules: []*typictx.Module{
		{ConfigPrefix: "APP", ConfigSpec: &config.Config{}},
		typserver.Module(),
		typpostgres.Module(),
	},

	Application: typictx.Application{
		StartFunc: func(s *app.Server, cfg *config.Config) error {
			log.Info("Start the application")
			return s.Start(cfg.Address)
		},
	},

	Release: typictx.Release{
		GoOS:   []string{"linux", "darwin"},
		GoArch: []string{"amd64"},
		Github: &typictx.Github{
			Owner:    "typical-go",
			RepoName: "typical-rest-server",
		},
	},

	DockerCompose: docker.NewCompose("3"),
}
