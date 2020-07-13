package typical

import (
	"github.com/typical-go/typical-go/pkg/typdocker"
	"github.com/typical-go/typical-go/pkg/typgo"
	"github.com/typical-go/typical-go/pkg/typmock"
	"github.com/typical-go/typical-rest-server/internal/app"
	"github.com/typical-go/typical-rest-server/internal/infra"
	"github.com/typical-go/typical-rest-server/pkg/dockerrx"
	"github.com/typical-go/typical-rest-server/pkg/pgcmd"
	"github.com/typical-go/typical-rest-server/pkg/rediscmd"
)

// Descriptor of Typical REST Server
// Build-Tool and Application will be generated based on this descriptor
var Descriptor = typgo.Descriptor{
	Name:        "typical-rest-server",
	Description: "Example of typical and scalable RESTful API Server for Go",
	Version:     "0.8.34",

	EntryPoint: app.Main,
	Layouts:    []string{"internal", "pkg"},

	Prebuild: typgo.Prebuilds{
		&typgo.DependencyInjection{},
		&typgo.ConfigManager{
			Configs: []*typgo.Configuration{
				{Name: "APP", Spec: &infra.App{}},
				{Name: "REDIS", Spec: &infra.Redis{}},
				{Name: "PG", Spec: &infra.Pg{}},
			},
		},
	},

	Compile: &typgo.StdCompile{},
	Run:     &typgo.StdRun{},
	Test:    &typgo.StdTest{},
	Clean:   &typgo.StdClean{},
	Release: &typgo.Github{Owner: "typical-go", RepoName: "typical-rest-server"},

	Utility: typgo.Utilities{
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
		&rediscmd.Utility{
			Name:        "redis",
			HostEnv:     "REDIS_HOST",
			PortEnv:     "REDIS_PORT",
			PasswordEnv: "REDIS_PASSWORD",
		},
		&typmock.Utility{},
		&typdocker.Utility{
			Version: typdocker.V3,
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
	},
}
