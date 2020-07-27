package main

import (
	"github.com/typical-go/typical-go/pkg/typannot"
	"github.com/typical-go/typical-go/pkg/typapp"
	"github.com/typical-go/typical-go/pkg/typdocker"
	"github.com/typical-go/typical-go/pkg/typgo"
	"github.com/typical-go/typical-go/pkg/typmock"
	"github.com/typical-go/typical-go/pkg/typrls"
	"github.com/typical-go/typical-rest-server/pkg/dockerrx"
	"github.com/typical-go/typical-rest-server/pkg/pgcmd"
)

var descriptor = typgo.Descriptor{
	Name:    "typical-rest-server",
	Version: "0.8.35",
	Layouts: []string{"internal", "pkg"},

	Cmds: []typgo.Cmd{

		&typgo.TestCmd{
			Action: &typgo.StdTest{},
		},

		&typgo.CompileCmd{
			Before: &typannot.Annotators{
				&typapp.CtorAnnotation{},
				&typapp.DtorAnnotation{},
				&typapp.CfgAnnotation{DotEnv: true},
			},
			Action: &typgo.StdCompile{},
		},

		&typgo.RunCmd{
			Before: typgo.BuildSysRuns{"compile"},
			Action: &typgo.StdRun{},
		},

		&typgo.CleanCmd{
			Action: &typgo.StdClean{},
		},

		&typmock.MockCmd{},

		&typdocker.DockerCmd{
			Composers: []typdocker.Composer{
				&dockerrx.PostgresWithEnv{
					Name:        "pg01",
					UserEnv:     "PG_USER",
					PasswordEnv: "PG_PASSWORD",
					PortEnv:     "PG_PORT",
				},
				&dockerrx.RedisWithEnv{
					Name:        "redis01",
					PasswordEnv: "REDIS_PASSWORD",
					PortEnv:     "REDIS_PORT",
				},
			},
		},

		&pgcmd.Utility{
			Name:         "pg",
			HostEnv:      "PG_HOST",
			PortEnv:      "PG_PORT",
			UserEnv:      "PG_USER",
			PasswordEnv:  "PG_PASSWORD",
			DBNameEnv:    "PG_DBNAME",
			MigrationSrc: "databases/pg/migration",
			SeedSrc:      "databases/pg/seed",
		},

		&typrls.ReleaseCmd{
			Before:     typgo.BuildSysRuns{"test", "compile"},
			Validation: typrls.DefaultValidation,
			Summary:    typrls.DefaultSummary,
			Releaser:   &typrls.Github{Owner: "typical-go", Repo: "typical-rest-server"},
		},
	},
}

func main() {
	typgo.Start(&descriptor)
}
