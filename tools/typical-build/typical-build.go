package main

import (
	"github.com/typical-go/typical-go/pkg/typapp"
	"github.com/typical-go/typical-go/pkg/typast"
	"github.com/typical-go/typical-go/pkg/typgo"
	"github.com/typical-go/typical-go/pkg/typmock"
	"github.com/typical-go/typical-go/pkg/typrls"
	"github.com/typical-go/typical-rest-server/pkg/dbtool"
	"github.com/typical-go/typical-rest-server/pkg/typcfg"
	"github.com/typical-go/typical-rest-server/pkg/typdocker"
)

var descriptor = typgo.Descriptor{
	ProjectName:    "typical-rest-server",
	ProjectVersion: "0.9.14",
	Environment:    typgo.DotEnv(".env"),

	Tasks: []typgo.Tasker{
		// annotate
		&typast.AnnotateProject{
			Sources: []string{"internal"},
			Annotators: []typast.Annotator{
				&typapp.CtorAnnotation{},
				// &typrepo.EntityAnnotation{},
				&typcfg.EnvconfigAnnotation{DotEnv: ".env", UsageDoc: "USAGE.md"},
			},
		},
		// test
		&typgo.GoTest{
			Includes: []string{"internal/app/**", "pkg/**"},
		},
		// compile
		&typgo.GoBuild{},
		// run
		&typgo.RunBinary{
			Before: typgo.TaskNames{"annotate", "build"},
		},
		// mock
		&typmock.GoMock{
			Sources: []string{"internal"},
		},
		// docker
		&typdocker.DockerTool{
			ComposeFiles: typdocker.ComposeFiles("deploy/docker"),
			EnvFile:      ".env",
		},
		// pg
		&dbtool.Postgres{
			Name:         "pg",
			DockerName:   "typical-rest-server-pg",
			EnvKeys:      dbtool.EnvKeysWithPrefix("PG"),
			MigrationSrc: "databases/pg/migration",
			SeedSrc:      "databases/pg/seed",
		},
		// mysql
		&dbtool.MySQL{
			Name:         "mysql",
			DockerName:   "typical-rest-server-mysql",
			EnvKeys:      dbtool.EnvKeysWithPrefix("MYSQL"),
			MigrationSrc: "databases/mysql/migration",
			SeedSrc:      "databases/mysql/seed",
		},
		// setup
		&typgo.Task{
			Name:  "setup",
			Usage: "setup dependency",
			Action: typgo.TaskNames{
				"pg drop", "pg create", "pg migrate", "pg seed",
				"mysql drop", "mysql create", "mysql migrate", "mysql seed",
			},
		},
		// release
		&typrls.ReleaseProject{
			Before: typgo.TaskNames{"test", "build"},
			// Releaser:  &typrls.CrossCompiler{Targets: []typrls.Target{"darwin/amd64", "linux/amd64"}},
			Publisher: &typrls.Github{Owner: "typical-go", Repo: "typical-rest-server"},
		},
	},
}

func main() {
	typgo.Start(&descriptor)
}
