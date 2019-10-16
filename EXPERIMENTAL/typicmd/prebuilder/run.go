package prebuilder

import (
	"os"

	log "github.com/sirupsen/logrus"

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
	log.Debug("Validate the context")
	fatalIfError(ctx.Validate())
	log.Debug("Prepare Environment File")
	typienv.PrepareEnvFile(ctx)
	prebuilder := prebuilder{}
	fatalIfError(prebuilder.Initiate(ctx))
	fatalIfError(prebuilder.Prebuild())
}

func fatalIfError(err error) {
	if err != nil {
		log.Fatal(err.Error())
	}
}
