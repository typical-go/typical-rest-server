package typitask

import (
	log "github.com/sirupsen/logrus"
	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/internal/bash"
	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/typictx"
	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/typienv"
)

// BuildBinary for typical application
func BuildBinary(ctx typictx.ActionContext) error {
	binaryName := typienv.Binary(ctx.Typical.BinaryNameOrDefault())
	mainPackage := typienv.AppMainPackage()

	log.Infof("Build the Binary for '%s' at '%s'", mainPackage, binaryName)
	bash.GoBuild(binaryName, mainPackage)

	return nil
}
