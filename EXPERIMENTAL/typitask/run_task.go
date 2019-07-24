package typitask

import (
	log "github.com/sirupsen/logrus"
	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/internal/bash"
	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/typictx"
)

// RunTest to run automate testing
func RunTest(ctx typictx.ActionContext) error {
	log.Info("Run the Test")
	bash.GoTest(ctx.Typical.AppModule.GetTestTargets())
	return nil
}
