package typicmd

import (
	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/bash"
	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/typictx"
	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/typienv"
)

// BuildBinary for typical application
func BuildBinary(ctx *typictx.ActionContext) error {
	binaryName := typienv.Binary(ctx.BinaryNameOrDefault())
	mainPackage := typienv.AppMainPackage()

	return bash.GoBuild(binaryName, mainPackage)
}
