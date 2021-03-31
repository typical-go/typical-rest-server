package main

import (
	"github.com/typical-go/typical-go/pkg/typapp"
	"github.com/typical-go/typical-go/pkg/typast"
	"github.com/typical-go/typical-go/pkg/typgo"
	"github.com/typical-go/typical-go/pkg/typmock"
	"github.com/typical-go/typical-go/pkg/typrls"
	"github.com/typical-go/typical-rest-server/pkg/typcfg"
	"github.com/typical-go/typical-rest-server/pkg/typrepo"
	"github.com/typical-go/typical-rest-server/pkg/typtool"
)

var descriptor = typgo.Descriptor{
	ProjectName:    "typical-rest-server",
	ProjectVersion: "0.9.15",
	Environment:    typgo.DotEnv(".env"),

	Tasks: []typgo.Tasker{
		// annotate
		&typast.AnnotateProject{
			Sources: []string{"internal"},
			Annotators: []typast.Annotator{
				&typapp.CtorAnnotation{},
				&typrepo.EntityAnnotation{},
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
		&typtool.Docker{
			ComposeFiles: typtool.ComposeFiles("deploy/docker"),
			EnvFile:      ".env",
		},
		// pg
		&typtool.Postgres{
			Name:         "pg",
			EnvKeys:      typtool.DBEnvKeysWithPrefix("PG"),
			MigrationSrc: "databases/pg/migration",
			SeedSrc:      "databases/pg/seed",
		},
		// mysql
		&typtool.MySQL{
			Name:         "mysql",
			EnvKeys:      typtool.DBEnvKeysWithPrefix("MYSQL"),
			MigrationSrc: "databases/mysql/migration",
			SeedSrc:      "databases/mysql/seed",
		},
		// mysql
		&typtool.Redis{
			Name:    "cache",
			EnvKeys: typtool.RedisEnvKeysWithPrefix("CACHE_REDIS"),
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
