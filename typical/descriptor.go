package typical

import (
	"github.com/typical-go/typical-go/pkg/typapp"
	"github.com/typical-go/typical-go/pkg/typbuildtool"
	"github.com/typical-go/typical-go/pkg/typcore"
	"github.com/typical-go/typical-go/pkg/typdocker"
	"github.com/typical-go/typical-rest-server/pkg/typpostgres"
	"github.com/typical-go/typical-rest-server/pkg/typredis"
	"github.com/typical-go/typical-rest-server/server"
)

// Descriptor of Typical REST Server
// Build-Tool and Application will be generated based on this descriptor
var Descriptor = typcore.Descriptor{
	Name:        "typical-rest-server",
	Description: "Example of typical and scalable RESTful API Server for Go",
	Version:     "0.8.27",

	App: typapp.EntryPoint(server.Main, "server").
		Imports(
			server.Configuration(),
			typredis.Module(),    // create/destroy redis conn + its configuration
			typpostgres.Module(), // create/destroy pg conn  + its configuration
		),

	BuildTool: typbuildtool.
		BuildSequences(
			typbuildtool.StandardBuild(),
			typbuildtool.Github("typical-go", "typical-rest-server"), // Create release to Github
		).
		Utilities(
			typpostgres.Utility(), // create db, drop, migrate, seed, console, etc.
			typredis.Utility(),    // redis console

			typdocker.Compose(
				typpostgres.DockerRecipeV3(),
				typredis.DockerRecipeV3(),
			),
		),
}
