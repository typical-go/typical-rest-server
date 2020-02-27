package typical

import (
	"github.com/typical-go/typical-go/pkg/typapp"
	"github.com/typical-go/typical-go/pkg/typbuildtool"
	"github.com/typical-go/typical-go/pkg/typcfg"
	"github.com/typical-go/typical-go/pkg/typcore"
	"github.com/typical-go/typical-go/pkg/typrls"
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
	Version: "0.8.20",

	// Configuration for this project
	// Both App and Build-Tool typically using the same configuration
	Configuration: typcfg.New().
		WithConfigure(
			server,
			redis,
			postgres,
		),

	// Detail of this application
	App: typapp.New(server).

		// Dependency is what are provided in dig service-locator
		// and what to be destroyed after application stop
		AppendDependency(
			redis,    // Create and destroy redis connection
			postgres, // Create and destroy postgres db connection
		).

		// Preparation before start the application
		AppendPreparer(
			redis,    // Ping to Redis Server
			postgres, // Ping to Postgres Database
		),

	// BuildTool responsible to basic build needs and custom dev task
	BuildTool: typbuildtool.New().

		// Additional command to be register in Build-Tool
		AppendCommander(
			postgres, // Postgres utilities like create, drop, migrate, seed, etc.
			redis,    // Redis utilities
			docker,   // Generate dockercompose and spin up docker
			readme,   // Generate readme based on README.tmpl
			rails,    // Experimental to generate code like RoR
		).

		// Setting to release the project
		WithRelease(
			// By default it will create distribution for Darwin and Linux
			typrls.New().
				WithPublisher(
					// Create release and upload file to Github
					typrls.GithubPublisher("typical-go", "typical-rest-server"),
				),
		),
}
