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
// Build-Tool and Application will be generated based on this descriptor
var Descriptor = typcore.Descriptor{

	// Common project information
	// Name and Package is mandatory
	Name:        "Typical REST Server",
	Description: "Example of typical and scalable RESTful API Server for Go",
	Version:     "0.8.17",
	Package:     "github.com/typical-go/typical-rest-server",

	App: typcore.NewApp(application).

		// Dependency is what are provided in dig service-locator
		// and what to be destroyed after application stop
		WithDependency(
			server,   // create and destroy http server
			redis,    // create and destroy redis connection
			postgres, // create and destroy postgres db connection
		).

		// Preparation before start the application
		WithPrepare(
			redis,    // Ping to Redis Server
			postgres, // Ping to Postgres Database
		),

	Build: typcore.NewBuild().

		// Additional command to be register in Build-Tool
		WithCommands(
			docker,
			readme, // generate readme based on README.tmpl
			postgres,
			redis,
			rails,
		).

		// Setting to release the project
		// By default it will create distribution for Darwin and Linux
		WithRelease(typrls.New().
			WithPublisher(
				// Create release and upload file to Github
				typrls.GithubPublisher("typical-go", "typical-rest-server"),
			),
		),

	// Configuration for this project
	// Both Build-Tool and Application typically using same configuration
	Configuration: typcore.NewConfiguration().
		WithConfigure(
			application,
			server,
			redis,
			postgres,
		),
}

// Modules that required for the project
var (
	application = app.New()
	readme      = typreadme.New()
	rails       = typrails.New()
	server      = typserver.New().WithDebug(true)
	redis       = typredis.New()
	postgres    = typpostgres.New().WithDBName("sample")

	docker = typdocker.New().WithComposers(
		postgres,
		redis,
	)
)
