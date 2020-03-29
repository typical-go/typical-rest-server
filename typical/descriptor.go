package typical

import (
	"github.com/typical-go/typical-go/pkg/typapp"
	"github.com/typical-go/typical-go/pkg/typbuildtool"
	"github.com/typical-go/typical-go/pkg/typcfg"
	"github.com/typical-go/typical-go/pkg/typcore"
	"github.com/typical-go/typical-go/pkg/typdocker"
	"github.com/typical-go/typical-go/pkg/typreadme"
	"github.com/typical-go/typical-rest-server/pkg/typpostgres"
	"github.com/typical-go/typical-rest-server/pkg/typredis"
	"github.com/typical-go/typical-rest-server/pkg/typserver"
	"github.com/typical-go/typical-rest-server/server"
)

// Descriptor of Typical REST Server
// Build-Tool and Application will be generated based on this descriptor
var Descriptor = typcore.Descriptor{
	Name:        "typical-rest-server",                                       // name of the project
	Description: "Example of typical and scalable RESTful API Server for Go", // description of the project
	Version:     "0.8.25",                                                    // version of the project

	// Detail of this application
	App: typapp.EntryPoint(server.Main, "server").
		WithModules(
			typserver.Module(),
			typredis.Module(),    // create and destroy redis connection
			typpostgres.Module(), // create and destroy postgres db connection
		),

	// BuildTool responsible to basic build needs and custom dev task
	BuildTool: typbuildtool.
		BuildSequences(
			typbuildtool.StandardBuild(),
			typbuildtool.Github("typical-go", "typical-rest-server"), // Create release to Github
		).
		WithUtilities(
			typpostgres.Utility(), // create database, drop, migrate, seed, etc.
			typredis.Utility(),    // redis console
			typreadme.Generator(), // generate readme based on README.tmpl

			// Generate dockercompose and spin up docker
			typdocker.Compose(
				typpostgres.DockerRecipeV3(),
				typredis.DockerRecipeV3(),
			),
		),

	// ConfigManager handle the configuration. Both App and Build-Tool typically using the same configuration
	ConfigManager: typcfg.
		Configures(
			server.Configuration(),
			typpostgres.Configuration(),
			typredis.Configuration(),
		),
}
