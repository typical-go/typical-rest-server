package typical

import (
	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/typictx"
	"github.com/typical-go/typical-rest-server/app"
	"github.com/typical-go/typical-rest-server/app/config"
	"github.com/typical-go/typical-rest-server/pkg/module/typpostgres"
	"github.com/typical-go/typical-rest-server/pkg/module/typredis"
	"github.com/typical-go/typical-rest-server/pkg/module/typserver"
	"github.com/urfave/cli"
)

// Context instance of Context
var Context = &typictx.Context{
	Root:        "github.com/typical-go/typical-rest-server",
	Name:        "Typical-RESTful-Server",
	Description: "Example of typical and scalable RESTful API Server for Go",
	Application: typictx.Application{
		StartFunc: app.Start,
		Config:    typictx.NewConfig("APP", &config.Config{}),
		Initiations: []interface{}{
			app.Middlewares,
			app.Routes,
		},
		Commands: []cli.Command{
			{Name: "route", Description: "Print available API Routes", Action: app.DryRun},
		},
	},
	Modules: []interface{}{
		typserver.Module(),
		typpostgres.Module(),
		typredis.Module(),
	},
	Release: typictx.Release{
		Version: "0.8.4",
		Targets: []string{
			"linux/amd64",
			"darwin/amd64",
		},
		Github: &typictx.Github{
			Owner:    "typical-go",
			RepoName: "typical-rest-server",
		},
		// Tagging: typictx.Tagging{
		// 	WithGitBranch:       true,
		// 	WithLatestGitCommit: true,
		// },
	},
}
