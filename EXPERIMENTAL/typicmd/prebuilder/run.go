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
	if os.Getenv(debugEnv) != "" {
		log.SetLevel(log.DebugLevel)
	}
	log.Debug("Preparing the context")
	fatalIfError(ctx.Preparing())
	log.Debug("Prepare Environment File")
	typienv.PrepareEnvFile(ctx)
	prebuilder := prebuilder{}
	fatalIfError(prebuilder.Initiate(ctx))
	report, err := prebuilder.Prebuild()
	fatalIfError(err)
	buildToolBinaryNotExist := !filekit.Exists(typienv.BuildTool.BinPath)
	log.Debugf("buildToolBinaryNotExist: %t", buildToolBinaryNotExist)
	prebuildUpdated := report.Updated()
	log.Debugf("prebuildUpdated: %t", prebuildUpdated)
	haveBuildArgs := haveBuildArg()
	log.Debugf("haveBuildArgs: %t", haveBuildArgs)
	if buildToolBinaryNotExist || prebuildUpdated || haveBuildArgs {
		log.Info("Build the build-tool")
		fatalIfError(bash.GoBuild(
			typienv.BuildTool.BinPath,
			typienv.BuildTool.SrcPath,
		))
	}
}

func haveBuildArg() bool {
	if len(os.Args) > 1 {
		return os.Args[1] == "1"
	}
	return false
}

func fatalIfError(err error) {
	if err != nil {
		log.Fatal(err.Error())
	}
}
