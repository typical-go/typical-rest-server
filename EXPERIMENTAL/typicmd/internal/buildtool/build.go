package buildtool

import (
	"os"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/bash"
	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/typictx"
	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/typienv"
)

func buildBinary(ctx *typictx.ActionContext) error {
	binaryName := typienv.App.Binary
	mainPackage := typienv.App.MainPkg
	return bash.GoBuild(binaryName, mainPackage)
}

func runBinary(ctx *typictx.ActionContext) error {
	if !ctx.Cli.Bool("no-build") {
		buildBinary(ctx)
	}
	binaryPath := typienv.App.Binary
	return bash.Run(binaryPath, []string(ctx.Cli.Args())...)
}

func runTesting(ctx *typictx.ActionContext) error {
	return bash.GoTest(ctx.TestTargets)

}

func generateMock(ctx *typictx.ActionContext) (err error) {
	err = bash.GoGet("github.com/golang/mock/mockgen")
	if err != nil {
		return
	}
	mockPkg := typienv.Mock
	if !ctx.Cli.Bool("no-delete") {
		log.Infof("Clean mock package '%s'", mockPkg)
		os.RemoveAll(mockPkg)
	}
	for _, mockTarget := range ctx.MockTargets {
		dest := mockPkg + "/" + mockTarget[strings.LastIndex(mockTarget, "/")+1:]
		err = bash.RunGoBin("mockgen",
			"-source", mockTarget,
			"-destination", dest,
			"-package", mockPkg)
	}
	return
}
