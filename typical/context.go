package typical

import (
	"github.com/typical-go/typical-go/pkg/typctx"
	"github.com/typical-go/typical-go/pkg/typrls"
	"github.com/typical-go/typical-rest-server/app"
	"github.com/typical-go/typical-rest-server/pkg/typdocker"
	"github.com/typical-go/typical-rest-server/pkg/typpostgres"
	"github.com/typical-go/typical-rest-server/pkg/typreadme"
	"github.com/typical-go/typical-rest-server/pkg/typredis"
	"github.com/typical-go/typical-rest-server/pkg/typserver"
)

// Context of project
var Context = &typctx.Context{
	Name:        "Typical-RESTful-Server",
	Description: "Example of typical and scalable RESTful API Server for Go",
	Package:     "github.com/typical-go/typical-rest-server",
	AppModule:   app.Module(),
	Modules: []interface{}{
		typdocker.Module(),
		typserver.Module(),
		typpostgres.Module(),
		typredis.Module(),
	},
	Releaser: typrls.Releaser{
		Version: "0.8.5",
		Targets: []typrls.ReleaseTarget{"linux/amd64", "darwin/amd64"},
		Publisher: &typrls.Github{
			Owner:    "typical-go",
			RepoName: "typical-rest-server",
		},
		// Tagging: typctx.Tagging{
		// 	WithGitBranch:       true,
		// 	WithLatestGitCommit: true,
		// },
	},
	ReadmeGenerator: typreadme.Generator{},
}
