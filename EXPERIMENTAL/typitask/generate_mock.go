package typitask

import (
	"os"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/internal/bash"
	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/typictx"
	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/typienv"
)

// GenerateMock to generate mock for interface as defined in mockTarget
func GenerateMock(ctx *typictx.ActionContext) (err error) {
	err = bash.GoGet("github.com/golang/mock/mockgen")
	if err != nil {
		return
	}

	mockPkg := typienv.Mock()

	if ctx.Cli.Bool("new") {
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
