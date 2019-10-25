package prebuilder

import (
	"os"

	"github.com/typical-go/typical-rest-server/pkg/utility/filekit"

	log "github.com/sirupsen/logrus"

	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/bash"
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

// Run the prebuilder
func Run(ctx *typictx.Context) {
	var prebuilder prebuilder
	var report report
	var err error
	if os.Getenv(debugEnv) != "" {
		log.SetLevel(log.DebugLevel)
	}
	log.Debug("Preparing the context")
	if err = ctx.Preparing(); err != nil {
		log.Fatal(err.Error())
	}
	log.Debug("Prepare Environment File")
	if err = typienv.PrepareEnvFile(ctx); err != nil {
		log.Fatal(err.Error())
	}
	log.Debug("Initiate prebuilder")
	if err := prebuilder.Initiate(ctx); err != nil {
		log.Fatal(err.Error())
	}
	log.Debug("Prebuilding")
	if report, err = prebuilder.Prebuild(); err != nil {
		log.Fatal(err.Error())
	}
	checker := buildToolChecker{
		BinaryNotExist:  !filekit.Exists(typienv.BuildTool.BinPath),
		PrebuildUpdated: report.Updated(),
		HaveBuildArgs:   haveBuildArg(),
	}
	if checker.Check() {
		log.Info("Build the build-tool")
		if err := bash.GoBuild(typienv.BuildTool.BinPath, typienv.BuildTool.SrcPath); err != nil {
			log.Fatal(err.Error())
		}
	}
}

func haveBuildArg() bool {
	if len(os.Args) > 1 {
		return os.Args[1] == "1"
	}
	return false
}
