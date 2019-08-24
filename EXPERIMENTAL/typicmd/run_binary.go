package typicmd

import (
	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/bash"
	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/typictx"
	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/typienv"
)

// RunBinary for run typical binary
func RunBinary(ctx *typictx.ActionContext) error {
	if !ctx.Cli.Bool("no-build") {
		BuildBinary(ctx)
	}

	binaryPath := typienv.Binary(ctx.BinaryNameOrDefault())
	return bash.Run(binaryPath, []string(ctx.Cli.Args())...)
}
