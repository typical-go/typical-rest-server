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
func GenerateMock(ctx typictx.ActionContext) error {
	bash.GoGet("github.com/golang/mock/mockgen")

	mockPkg := typienv.Mock()

	if ctx.Cli.Bool("new") {
		log.Infof("Clean mock package '%s'", mockPkg)
		os.RemoveAll(mockPkg)
	}

	for _, mockTarget := range ctx.Typical.MockTargets {
		dest := mockPkg + "/" + mockTarget[strings.LastIndex(mockTarget, "/")+1:]

		log.Infof("Generate mock for '%s' at '%s'", mockTarget, dest)
		bash.RunGoBin("mockgen",
			"-source", mockTarget,
			"-destination", dest,
			"-package", mockPkg)
	}
	return nil
}
