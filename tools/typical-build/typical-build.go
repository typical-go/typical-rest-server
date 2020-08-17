package main

import (
	"github.com/typical-go/typical-go/pkg/typannot"
	"github.com/typical-go/typical-go/pkg/typapp"
	"github.com/typical-go/typical-go/pkg/typgo"
	"github.com/typical-go/typical-go/pkg/typmock"
	"github.com/typical-go/typical-go/pkg/typrls"
	"github.com/typical-go/typical-rest-server/pkg/typcfg"
	"github.com/typical-go/typical-rest-server/pkg/typdocker"
	"github.com/typical-go/typical-rest-server/tools/typical-build/pg"
)

var descriptor = typgo.Descriptor{
	ProjectName:    "typical-rest-server",
	ProjectVersion: "0.8.37",
	ProjectLayouts: []string{"internal", "pkg"},

	Cmds: []typgo.Cmd{
		// clean
		&typgo.CleanProject{},
		// test
		&typgo.TestProject{},
		// compile
		&typgo.CompileProject{},
		// annotate
		&typannot.AnnotateCmd{
			Annotators: []typannot.Annotator{
				&typapp.CtorAnnotation{},
				&typapp.DtorAnnotation{},
				&typcfg.AppCfgAnnotation{
					DotEnv: true, // generate .env file reflected to config
				},
			},
		},
		// run
		&typgo.RunCmd{
			Before: typgo.BuildCmdRuns{"annotate", "compile"},
			Action: &typgo.RunProject{},
		},
		// mock
		&typmock.MockCmd{},
		// docker
		&typdocker.DockerCmd{},
		// pg
		&pg.Utility{},
		// reset
		&typgo.Command{
			Name:  "reset",
			Usage: "reset the project locally (postgres/etc)",
			Action: typgo.BuildCmdRuns{
				"pg.drop",
				"pg.create",
				"pg.migrate",
				"pg.seed",
			},
		},
		// release
		&typrls.ReleaseCmd{
			Before: typgo.BuildCmdRuns{"test", "compile"},
			Action: &typrls.ReleaseProject{
				Validator: typrls.DefaultValidator,
				Releaser:  &typrls.Github{Owner: "typical-go", Repo: "typical-rest-server"},
			},
		},
	},
}

func main() {
	typgo.Start(&descriptor)
}
