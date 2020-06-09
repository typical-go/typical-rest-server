package typical

import (
	"github.com/typical-go/typical-go/pkg/typdocker"
	"github.com/typical-go/typical-go/pkg/typgo"
	"github.com/typical-go/typical-go/pkg/typmock"
	"github.com/typical-go/typical-rest-server/internal/app"
	"github.com/typical-go/typical-rest-server/internal/app/infra"
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
	Layouts:    []string{"internal", "pkg"},

	Prebuild: typgo.Prebuilds{
		&typgo.DependencyInjection{},
		&typgo.ConfigManager{
			Configs: []*typgo.Configuration{
				{Name: "APP", Spec: &infra.App{Debug: true}},
				{Name: "REDIS", Spec: &infra.Redis{}},
				{Name: "PG", Spec: &typpg.Config{DBName: "MyLibrary"}},
			},
		},
	},

	Compile: &typgo.StdCompile{},
	Run:     &typgo.StdRun{},
	Test:    &typgo.StdTest{},
	Clean:   &typgo.StdClean{},
	Release: &typgo.Github{Owner: "typical-go", RepoName: "typical-rest-server"},

	Utility: typgo.Utilities{
		&typpg.Utility{
			Name:         "pg",
			MigrationSrc: "scripts/db/migration",
			SeedSrc:      "scripts/db/seed",
			ConfigName:   "PG",
		},
		typgo.NewUtility(redisUtil),
		&typmock.Utility{},
		&typdocker.Utility{
			Composers: []typdocker.Composer{
				&pgDocker{},
				&redisDocker{},
			},
		},
	},
}
