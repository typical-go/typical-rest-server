package prebuilder

import (
	"os"

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
	ctxPath  = "typical/context.go"
	debugEnv = "PREBUILDER_DEBUG"
)

// Prebuild process
func Prebuild(ctx *typictx.Context) {
	if os.Getenv(debugEnv) != "" {
		log.SetLevel(log.DebugLevel)
	}
	log.Debug("Validate the context")
	fatalIfError(ctx.Validate())
	root := typienv.AppName
	log.Debug("Scan project to get package and filenames")
	packages, filenames, err := scanProject(root)
	fatalIfError(err)
	log.Debug("Walk the project to get annotated or metadata")
	projFiles, err := walker.WalkProject(filenames)
	fatalIfError(err)
	log.Debug("Walk the context file")
	ctxFile, err := walker.WalkContext(ctxPath)
	fatalIfError(err)
	log.Debug("Prepare Environment File")
	typienv.PrepareEnvFile(ctx)
	prebuilder := prebuilder{
		Context:      ctx,
		Filenames:    filenames,
		Packages:     packages,
		ProjectFiles: projFiles,
		ContextFile:  ctxFile,
	}
	if prebuilder.checkTestTargets() {
		fatalIfError(prebuilder.generateTestTargets())
	}
	if prebuilder.checkAnnotated() {
		fatalIfError(prebuilder.generateAnnotated())
	}
	if prebuilder.checkConfiguration() {
		fatalIfError(prebuilder.generateConfiguration())
	}
}

func fatalIfError(err error) {
	if err != nil {
		log.Fatal(err.Error())
	}
}
