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
	"github.com/typical-go/typical-rest-server/server"
)

// Modules that required for the project
var (
	serverApp = server.New().WithDebug(true)
	postgres  = typpostgres.New().WithDBName("sample")
	redis     = typredis.New()
	// rails     = typrails.New()
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
	Version: "0.8.23",

	// Detail of this application
	App: typapp.
		Create(serverApp).
		WithModules(
			redis,    // Create and destroy redis connection
			postgres, // Create and destroy postgres db connection
		),

	// Configuration for this project
	// Both App and Build-Tool typically using the same configuration
	ConfigManager: typcfg.
		Create(
			serverApp,
			redis,
			postgres,
		),

	// BuildTool responsible to basic build needs and custom dev task
	BuildTool: typbuildtool.
		Create(
			typbuildtool.StandardBuild(),
			typbuildtool.Github("typical-go", "typical-rest-server"), // Create release and upload file to Github
		).
		WithCommanders(
			// Postgres utilities like create, drop, migrate, seed, etc.
			postgres,

			// Generate dockercompose and spin up docker
			typdocker.
				Create().
				WithComposers(
					postgres,
					redis,
				),

			// Generate readme based on README.tmpl
			typreadme.Create(),

			// redis,    // Redis utilities
			// rails,    // Experimental to generate code like RoR
		),
}
