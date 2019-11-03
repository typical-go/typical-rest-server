package typical

import (
	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/typictx"
	"github.com/typical-go/typical-rest-server/app"
	"github.com/typical-go/typical-rest-server/pkg/typpostgres"
	"github.com/typical-go/typical-rest-server/pkg/typredis"
	"github.com/typical-go/typical-rest-server/pkg/typserver"
)

// Context instance of Context
var Context = &typictx.Context{
	Name:        "Typical-RESTful-Server",
	Description: "Example of typical and scalable RESTful API Server for Go",
	Root:        "github.com/typical-go/typical-rest-server",
	AppModule:   app.Module(),
	Modules: []interface{}{
		typserver.Module(),
		typpostgres.Module(),
		typredis.Module(),
	},
	Release: typictx.Release{
		Version: "0.8.4",
		Targets: []string{"linux/amd64", "darwin/amd64"},
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
