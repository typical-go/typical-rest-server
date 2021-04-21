package main

import (
	"github.com/typical-go/typical-go/pkg/typapp"
	"github.com/typical-go/typical-go/pkg/typast"
	"github.com/typical-go/typical-go/pkg/typgo"
	"github.com/typical-go/typical-go/pkg/typmock"
	"github.com/typical-go/typical-go/pkg/typrls"
	"github.com/typical-go/typical-rest-server/pkg/dbrepo"
	"github.com/typical-go/typical-rest-server/pkg/typcfg"
	"github.com/typical-go/typical-rest-server/pkg/typtool"
)

var descriptor = typgo.Descriptor{
	ProjectName:    "typical-rest-server",
	ProjectVersion: "0.9.16",
	Environment:    typgo.DotEnv(".env"),

	Tasks: []typgo.Tasker{
		// annotate
		&typast.AnnotateProject{
			Sources: []string{"internal"},
			Annotators: []typast.Annotator{
				&typapp.CtorAnnotation{},
				&dbrepo.DBRepoAnnot{},
				&typcfg.EnvconfigAnnot{GenDotEnv: ".env", GenDoc: "USAGE.md"},
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
			MigrationSrc: "database/pg/migration",
			SeedSrc:      "database/pg/seed",
		},
		// mysql
		&typtool.MySQL{
			Name:         "mysql",
			EnvKeys:      typtool.DBEnvKeysWithPrefix("MYSQL"),
			MigrationSrc: "database/mysql/migration",
			SeedSrc:      "database/mysql/seed",
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
