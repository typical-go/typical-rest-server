package typical

import (
	"github.com/typical-go/typical-go/pkg/typdocker"
	"github.com/typical-go/typical-go/pkg/typgo"
	"github.com/typical-go/typical-go/pkg/typmock"
	"github.com/typical-go/typical-rest-server/pkg/typpostgres"
	"github.com/typical-go/typical-rest-server/pkg/typredis"
	"github.com/typical-go/typical-rest-server/server"
)

// Descriptor of Typical REST Server
// Build-Tool and Application will be generated based on this descriptor
var Descriptor = typgo.Descriptor{

	Name:        "typical-rest-server",
	Description: "Example of typical and scalable RESTful API Server for Go",
	Version:     "0.8.28",

	EntryPoint: server.Main,

	Layouts: []string{"server", "pkg"},

	Configurer: typgo.Configurers{
		typredis.Configuration(),
		typpostgres.Configuration(nil),
		server.Configuration(),
	},

	Build: typgo.Builds{
		&typgo.StdBuild{},
		&typgo.Github{
			Owner:    "typical-go",
			RepoName: "typical-rest-server",
		},
	},

	Utility: typgo.Utilities{
		typpostgres.Utility(nil), // create db, drop, migrate, seed, console, etc.
		typredis.Utility(),       // redis console
		typmock.Utility(),

		typdocker.Compose(
			typpostgres.DockerRecipeV3(nil),
			typredis.DockerRecipeV3(),
		),
	},
}
