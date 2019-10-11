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

const (
	ctxPath = "typical/context.go"
)

// Prebuild process
func Prebuild(ctx *typictx.Context) (err error) {
	log.Debug("Validate the context")
	err = ctx.Validate()
	if err != nil {
		return
	}
	root := typienv.AppName
	log.Debug("Scan project to get package and filenames")
	packages, filenames, err := scanProject(root)
	if err != nil {
		return
	}
	log.Debug("Walk the project to get annotated or metadata")
	projFiles, err := walker.WalkProject(filenames)
	if err != nil {
		return
	}
	log.Debug("Walk the context file")
	ctxFile, err := walker.WalkContext(ctxPath)
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
	log.Debug("Prepare Environment File")
	typienv.PrepareEnvFile(ctx)
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
