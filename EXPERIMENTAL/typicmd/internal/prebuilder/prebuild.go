package prebuilder

import (
	"github.com/typical-go/runn"
	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/typicmd/internal/prebuilder/walker"
	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/typictx"
	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/typienv"
)

var (
	app        = typienv.App.SrcPath
	buildTool  = typienv.BuildTool.SrcPath
	dependency = typienv.Dependency.SrcPath
)

// PreBuild process to build the typical project
func PreBuild(ctx *typictx.Context) (err error) {
	root := typienv.AppName
	packages, filenames, _ := projectFiles(root)
	projFiles, err := walker.WalkProject(filenames)
	if err != nil {
		return
	}
	ctxFile, err := walker.WalkContext("typical/context.go")
	if err != nil {
		return
	}
	prebuilder := PreBuilder{
		Context:      ctx,
		Filenames:    filenames,
		Packages:     packages,
		ProjectFiles: projFiles,
		ContextFile:  ctxFile,
	}
	return runn.Execute(
		typienv.WriteEnvIfNotExist(ctx),
		prebuilder.TestTargets(),
		prebuilder.Annotated(),
		prebuilder.Configuration(),
	)
}
