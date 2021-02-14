package main

import (
	"github.com/typical-go/typical-go/pkg/typapp"
	"github.com/typical-go/typical-go/pkg/typast"
	"github.com/typical-go/typical-go/pkg/typgo"
	"github.com/typical-go/typical-go/pkg/typmock"
	"github.com/typical-go/typical-go/pkg/typrls"
	"github.com/typical-go/typical-rest-server/pkg/dbtool"
	"github.com/typical-go/typical-rest-server/pkg/dbtool/mysqltool"
	"github.com/typical-go/typical-rest-server/pkg/dbtool/pgtool"
	"github.com/typical-go/typical-rest-server/pkg/typcfg"
	"github.com/typical-go/typical-rest-server/pkg/typdocker"
	"github.com/typical-go/typical-rest-server/pkg/typrepo"
)

var descriptor = typgo.Descriptor{
	ProjectName:    "typical-rest-server",
	ProjectVersion: "0.9.13",
	EnvLoader:      typgo.DotEnv(".env"),

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
			Excludes: []string{"internal/app/model"},
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
		&typdocker.DockerTool{},
		// pg
		&pgtool.PgTool{
			Name: "pg",
			Config: dbtool.Config{ // based on docker-compose.yml
				DBName: "dbname",
				DBUser: "dbuser",
				DBPass: "dbpass",
				Host:   "localhost",
				Port:   "5432",
			},
			DockerName:   "typical-rest-server_pg01_1",
			MigrationSrc: "databases/pg/migration",
			SeedSrc:      "databases/pg/seed",
		},
		// mysql
		&mysqltool.MySQLTool{
			Name: "mysql",
			Config: dbtool.Config{ // based on docker-compose.yml
				DBName: "dbname",
				DBUser: "dbuser",
				DBPass: "dbpass",
				Host:   "localhost",
				Port:   "3306",
			},
			DockerName:   "typical-rest-server_mysql01_1",
			MigrationSrc: "databases/mysql/migration",
			SeedSrc:      "databases/mysql/seed",
		},
		// reset
		&typgo.Task{
			Name:  "reset",
			Usage: "reset the project locally (postgres/etc)",
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
