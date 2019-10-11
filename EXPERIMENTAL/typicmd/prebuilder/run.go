package prebuilder

import (
	log "github.com/sirupsen/logrus"

	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/typicmd/prebuilder/walker"
	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/typictx"
	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/typienv"
)

var (
	app        = typienv.App.SrcPath
	buildTool  = typienv.BuildTool.SrcPath
	dependency = typienv.Dependency.SrcPath
)

// Prebuild process
func Prebuild(ctx *typictx.Context) (err error) {
	log.Debug("Validate the context")
	err = ctx.Validate()
	if err != nil {
		return
	}
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
	typienv.WriteEnvIfNotExist(ctx)
	err = prebuilder.TestTargets()
	if err != nil {
		return
	}
	err = prebuilder.Annotated()
	if err != nil {
		return
	}
	return prebuilder.Configuration()

}
