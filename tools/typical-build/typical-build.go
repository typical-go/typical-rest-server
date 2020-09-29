package main

import (
	"github.com/typical-go/typical-go/pkg/typapp"
	"github.com/typical-go/typical-go/pkg/typast"
	"github.com/typical-go/typical-go/pkg/typgo"
	"github.com/typical-go/typical-go/pkg/typmock"
	"github.com/typical-go/typical-go/pkg/typrls"
	"github.com/typical-go/typical-rest-server/internal/generated/typical"
	"github.com/typical-go/typical-rest-server/pkg/typcfg"
	"github.com/typical-go/typical-rest-server/pkg/typdocker"
	"github.com/typical-go/typical-rest-server/tools/typical-build/mysql"
	"github.com/typical-go/typical-rest-server/tools/typical-build/pg"
)

var descriptor = typgo.Descriptor{
	ProjectName:    "typical-rest-server",
	ProjectVersion: "0.9.3",
	ProjectLayouts: []string{"internal", "pkg"},

	Cmds: []typgo.Cmd{
		// test
		&typgo.TestProject{},
		// compile
		&typgo.CompileProject{},
		// annotate
		&typast.AnnotateProject{
			Destination: "internal/generated/typical",
			Annotators: []typast.Annotator{
				&typapp.CtorAnnotation{},
				&typapp.DtorAnnotation{},
				&typcfg.EnvconfigAnnotation{
					DotEnv:   ".env",     // generate .env file
					UsageDoc: "USAGE.md", // generate USAGE.md
				},
			},
		},
		// run
		&typgo.RunProject{
			Before: typgo.BuildCmdRuns{"annotate", "compile"},
		},
		// mock
		&typmock.MockCmd{},
		// docker
		&typdocker.DockerCmd{},
		// pg
		&pg.Command{
			Name:         "pg",
			ConfigFn:     typical.LoadPostgresCfg,
			DockerName:   "typical-rest-server_pg01_1",
			MigrationSrc: "file://databases/librarydb/migration",
			SeedSrc:      "databases/librarydb/seed",
		},
		// mysql
		&mysql.Command{
			Name:         "mysql",
			ConfigFn:     typical.LoadMySQLCfg,
			DockerName:   "typical-rest-server_mysql01_1",
			MigrationSrc: "file://databases/albumdb/migration",
			SeedSrc:      "databases/albumdb/seed",
		},
		// reset
		&typgo.Command{
			Name:  "reset",
			Usage: "reset the project locally (postgres/etc)",
			Action: typgo.BuildCmdRuns{
				"pg.drop", "pg.create", "pg.migrate", "pg.seed",
				"mysql.drop", "mysql.create", "mysql.migrate", "mysql.seed",
			},
		},
		// release
		&typrls.ReleaseProject{
			Before: typgo.BuildCmdRuns{"test", "compile"},
			// Releaser:  &typrls.CrossCompiler{Targets: []typrls.Target{"darwin/amd64", "linux/amd64"}},
			Publisher: &typrls.Github{Owner: "typical-go", Repo: "typical-rest-server"},
		},
	},
}

func main() {
	typgo.Start(&descriptor)
}
