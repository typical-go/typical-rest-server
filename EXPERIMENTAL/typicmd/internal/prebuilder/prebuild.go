package prebuilder

import (
	log "github.com/sirupsen/logrus"
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
	log.Info("Prebuilding...")
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
