package main

import (
	"github.com/typical-go/typical-go/pkg/typannot"
	"github.com/typical-go/typical-go/pkg/typapp"
	"github.com/typical-go/typical-go/pkg/typgo"
	"github.com/typical-go/typical-go/pkg/typmock"
	"github.com/typical-go/typical-go/pkg/typrls"
	"github.com/typical-go/typical-rest-server/pkg/typdocker"
	"github.com/typical-go/typical-rest-server/pkg/typrest"
	"github.com/typical-go/typical-rest-server/tools/typical-build/pg"
)

var descriptor = typgo.Descriptor{
	ProjectName:    "typical-rest-server",
	ProjectVersion: "0.8.37",
	ProjectLayouts: []string{"internal", "pkg"},

	Cmds: []typgo.Cmd{
		// test
		&typgo.TestCmd{Action: &typgo.StdTest{}},
		// compile
		&typgo.CompileCmd{Action: &typgo.StdCompile{}},
		// clean
		&typgo.CleanCmd{Action: &typgo.StdClean{}},
		// annotate
		&typannot.AnnotateCmd{
			Annotators: []typannot.Annotator{
				&typapp.CtorAnnotation{},
				&typapp.DtorAnnotation{},
				&typrest.AppCfgAnnotation{DotEnv: true},
			},
		},
		// run
		&typgo.RunCmd{
			Before: typgo.BuildSysRuns{"annotate", "compile"},
			Action: &typgo.StdRun{},
		},
		// mock
		&typmock.MockCmd{},
		// docker
		&typdocker.DockerCmd{},
		// pg
		&pg.Cmd{},

		// release
		&typrls.ReleaseCmd{
			Before:     typgo.BuildSysRuns{"test", "compile"},
			Validation: typrls.DefaultValidation,
			Summary:    typrls.DefaultSummary,
			Tag:        typrls.DefaultTag,
			Releaser:   &typrls.Github{Owner: "typical-go", Repo: "typical-rest-server"},
		},
	},
}

func main() {
	typgo.Start(&descriptor)
}
