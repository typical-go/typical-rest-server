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

// Modules that required for the project
var (
	postgres = typpostgres.New()
	redis    = typredis.New()
)

// Descriptor of Typical REST Server
// Build-Tool and Application will be generated based on this descriptor
var Descriptor = typcore.Descriptor{

	// Name of the project
	// It should be a characters with/without underscore or dash.
	Name: "typical-rest-server",

	// Description of the project
	Description: "Example of typical and scalable RESTful API Server for Go",

	// Version of the project
	Version: "0.8.24",

	// Detail of this application
	App: typapp.EntryPoint(server.Main, "server").
		WithModules(
			typserver.Module(),
			redis,    // Create and destroy redis connection
			postgres, // Create and destroy postgres db connection
		),

	// BuildTool responsible to basic build needs and custom dev task
	BuildTool: typbuildtool.
		BuildSequences(
			typbuildtool.StandardBuild(),
			typbuildtool.Github("typical-go", "typical-rest-server"), // Create release and upload file to Github
		).
		WithTasks(
			// Postgres utilities like create, drop, migrate, seed, etc.
			postgres,

			// Generate dockercompose and spin up docker
			typdocker.Compose(
				postgres,
				redis,
			),

			// Generate readme based on README.tmpl
			typreadme.Create(),
		),

	// Configuration for this project
	// Both App and Build-Tool typically using the same configuration
	ConfigManager: typcfg.
		Configures(
			server.Configuration(),
			typpostgres.Configuration(),
			typredis.Configuration(),
		),
}
