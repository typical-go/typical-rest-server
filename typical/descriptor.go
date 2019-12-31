package typical

import (
	"github.com/typical-go/typical-go/pkg/typcore"
	"github.com/typical-go/typical-go/pkg/typrls"
	"github.com/typical-go/typical-rest-server/app"
	"github.com/typical-go/typical-rest-server/pkg/typdocker"
	"github.com/typical-go/typical-rest-server/pkg/typpostgres"
	"github.com/typical-go/typical-rest-server/pkg/typrails"
	"github.com/typical-go/typical-rest-server/pkg/typreadme"
	"github.com/typical-go/typical-rest-server/pkg/typredis"
	"github.com/typical-go/typical-rest-server/pkg/typserver"
)

// Descriptor of Typical REST Server
var Descriptor = &typcore.ProjectDescriptor{
	Name:        "Typical REST Server",
	Description: "Example of typical and scalable RESTful API Server for Go",
	Version:     "0.8.14",
	Package:     "github.com/typical-go/typical-rest-server",

	AppModule: app.Module(),

	Modules: []interface{}{
		// General
		typdocker.Module(),
		typreadme.Module(),
		typrails.Module(),

		// HTTP Server
		typserver.Module(),

		// Redis
		typredis.New(),

		// Database
		typpostgres.New().WithDBName("sample"),
	},

	Releaser: typrls.New().WithPublisher(
		typrls.GithubPublisher("typical-go", "typical-rest-server"),
	),
}
