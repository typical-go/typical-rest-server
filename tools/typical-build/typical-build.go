package main

import (
	"log"

	"github.com/typical-go/typical-go/pkg/typapp"
	"github.com/typical-go/typical-go/pkg/typast"
	"github.com/typical-go/typical-go/pkg/typgo"
	"github.com/typical-go/typical-go/pkg/typmock"
	"github.com/typical-go/typical-go/pkg/typrls"
	"github.com/typical-go/typical-rest-server/internal/generated/config"
	"github.com/typical-go/typical-rest-server/pkg/dbtool"
	"github.com/typical-go/typical-rest-server/pkg/dbtool/mysqltool"
	"github.com/typical-go/typical-rest-server/pkg/dbtool/pgtool"
	"github.com/typical-go/typical-rest-server/pkg/typcfg"
	"github.com/typical-go/typical-rest-server/pkg/typdocker"
	"github.com/typical-go/typical-rest-server/pkg/typrepo"
)

var descriptor = typgo.Descriptor{
	ProjectName:    "typical-rest-server",
	ProjectVersion: "0.9.7",
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
				&typrepo.EntityAnnotation{},
				&typcfg.EnvconfigAnnotation{DotEnv: ".env", UsageDoc: "USAGE.md"},
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
		&pgtool.PgTool{
			Name: "pg",
			ConfigFn: func() dbtool.Configurer {
				cfg, err := config.LoadPgDatabaseCfg()
				if err != nil {
					log.Fatal(err)
				}
				return cfg
			},
			DockerName:   "typical-rest-server_pg01_1",
			MigrationSrc: "file://databases/postgresdb/migration",
			SeedSrc:      "databases/postgresdb/seed",
		},
		// mysql
		&mysqltool.MySQLTool{
			Name: "mysql",
			ConfigFn: func() dbtool.Configurer {
				cfg, err := config.LoadMysqlDatabaseCfg()
				if err != nil {
					log.Fatal(err)
				}
				return cfg
			},
			DockerName:   "typical-rest-server_mysql01_1",
			MigrationSrc: "file://databases/mysqldb/migration",
			SeedSrc:      "databases/mysqldb/seed",
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
