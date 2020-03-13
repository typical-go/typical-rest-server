package typical

import (
	"github.com/typical-go/typical-go/pkg/typapp"
	"github.com/typical-go/typical-go/pkg/typbuildtool"
	"github.com/typical-go/typical-go/pkg/typcfg"
	"github.com/typical-go/typical-go/pkg/typcore"
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
	Version: "0.8.22",

	// Detail of this application
	App: typapp.New(serverApp).
		AppendModule(
			redis,    // Create and destroy redis connection
			postgres, // Create and destroy postgres db connection
		),

	// Configuration for this project
	// Both App and Build-Tool typically using the same configuration
	ConfigManager: typcfg.New().
		WithConfigurers(
			serverApp,
			redis,
			postgres,
		),

	// BuildTool responsible to basic build needs and custom dev task
	BuildTool: typbuildtool.New().
		AppendCommanders(
			postgres, // Postgres utilities like create, drop, migrate, seed, etc.
			redis,    // Redis utilities
			docker,   // Generate dockercompose and spin up docker
			readme,   // Generate readme based on README.tmpl
			rails,    // Experimental to generate code like RoR
		).
		WithPublishers(
			typbuildtool.NewGithub("typical-go", "typical-rest-server"), // Create release and upload file to Github
		),
}
