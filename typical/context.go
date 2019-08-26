package typical

import (
	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/typictx"
	"github.com/typical-go/typical-rest-server/app"
	"github.com/typical-go/typical-rest-server/app/config"
	"github.com/typical-go/typical-rest-server/pkg/module/typpostgres"
	"github.com/typical-go/typical-rest-server/pkg/module/typserver"
)

// Context instance of Context
var Context = &typictx.Context{
	Name:        "Typical-RESTful-Server",
	Description: "Example of typical and scalable RESTful API Server for Go",
	Application: typictx.Application{
		StartFunc: app.Start,
		Config: typictx.Config{
			Prefix: "APP",
			Spec:   &config.Config{},
		},
		Initiations: []interface{}{
			app.Middlewares,
			app.Routes,
		},
	},
	Modules: []*typictx.Module{
		typserver.Module(),
		typpostgres.Module(),
	},
	Release: typictx.Release{
		Version: "0.7.0",
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
