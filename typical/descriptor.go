package typical

import (
	"github.com/typical-go/typical-go/pkg/typdocker"
	"github.com/typical-go/typical-go/pkg/typgo"
	"github.com/typical-go/typical-go/pkg/typmock"
	"github.com/typical-go/typical-rest-server/internal/app"
	"github.com/typical-go/typical-rest-server/internal/app/config"
	"github.com/typical-go/typical-rest-server/pkg/typpg"
)

var (
	mainDB = typpg.Init(&typpg.Settings{
		DBName:     "MyLibrary",
		DockerName: "pg01",
	})
)

// Descriptor of Typical REST Server
// Build-Tool and Application will be generated based on this descriptor
var Descriptor = typgo.Descriptor{

	Name:        "typical-rest-server",
	Description: "Example of typical and scalable RESTful API Server for Go",
	Version:     "0.8.31",

	EntryPoint: app.Main,

	Layouts: []string{
		"internal",
		"pkg",
	},

	Configurer: typgo.Configurers{
		&typgo.Configuration{Name: "APP", Spec: &config.Config{Debug: true}},
		&typgo.Configuration{Name: "REDIS", Spec: &config.Redis{}},
		typpg.Configuration(mainDB),
	},

	Build: typgo.Builds{
		&typgo.StdBuild{},
		&typgo.Github{
			Owner:    "typical-go",
			RepoName: "typical-rest-server",
		},
	},

	Utility: typgo.Utilities{
		typpg.Utility(mainDB), // create db, drop, migrate, seed, console, etc.
		typgo.NewUtility(redisUtil),
		typmock.Utility(),

		typdocker.Compose(
			typpg.DockerRecipeV3(mainDB),
			redisDocker(),
		),
	},
}
