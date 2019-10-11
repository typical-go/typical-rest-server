package buildtool

import (
	"os"
	"path/filepath"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/bash"
	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/typictx"
	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/typienv"
)

func buildBinary(ctx *typictx.ActionContext) error {
	binaryName := typienv.App.BinPath
	mainPackage := typienv.App.SrcPath
	return bash.GoBuild(binaryName, mainPackage)
}

func cleanProject(ctx *typictx.ActionContext) (err error) {
	err = os.RemoveAll(typienv.Bin)
	if err != nil {
		return
	}
	return filepath.Walk(typienv.Dependency.SrcPath, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			return os.Remove(path)
		}
		return nil
	})
}

func runBinary(ctx *typictx.ActionContext) error {
	if !ctx.Cli.Bool("no-build") {
		buildBinary(ctx)
	}
	binaryPath := typienv.App.BinPath
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
