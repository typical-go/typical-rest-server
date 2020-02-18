package typical

import (
	"github.com/typical-go/typical-go/pkg/typapp"
	"github.com/typical-go/typical-go/pkg/typbuildtool"
	"github.com/typical-go/typical-go/pkg/typbuildtool/stdrls"
	"github.com/typical-go/typical-go/pkg/typcfg"
	"github.com/typical-go/typical-go/pkg/typcore"
)

// Descriptor of Typical REST Server
// Build-Tool and Application will be generated based on this descriptor
var Descriptor = typcore.Descriptor{

	// Name is optional with default value is same with project folder
	// It should be a characters with/without underscore or dash.
	// Name: "typical-rest",

	// Description of the project
	Description: "Example of typical and scalable RESTful API Server for Go",

	// Version of the project
	Version: "0.8.19",

	// Package must be same with go.mod file
	Package: "github.com/typical-go/typical-rest-server",

	App: typapp.New(application).

		// Dependency is what are provided in dig service-locator
		// and what to be destroyed after application stop
		AppendDependency(
			server,   // create and destroy http server
			redis,    // create and destroy redis connection
			postgres, // create and destroy postgres db connection
		).

		// Preparation before start the application
		AppendPreparer(
			redis,    // Ping to Redis Server
			postgres, // Ping to Postgres Database
		),

	BuildTool: typbuildtool.New().

		// Additional command to be register in Build-Tool
		AppendCommander(
			docker,
			readme, // generate readme based on README.tmpl
			postgres,
			redis,
			rails,
		).

		// Setting to release the project
		// By default it will create distribution for Darwin and Linux
		WithRelease(stdrls.New().
			WithPublisher(
				// Create release and upload file to Github
				stdrls.GithubPublisher("typical-go", "typical-rest-server"),
			),
		),

	// Configuration for this project
	// Both Build-Tool and Application typically using same configuration
	Configuration: typcfg.New().
		WithConfigure(
			application,
			server,
			redis,
			postgres,
		),
}
