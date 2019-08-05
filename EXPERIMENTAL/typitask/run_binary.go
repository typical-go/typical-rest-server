package typitask

import (
	log "github.com/sirupsen/logrus"
	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/internal/bash"
	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/typictx"
	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/typienv"
)

// RunBinary for run typical binary
func RunBinary(ctx typictx.ActionContext) error {
	if !ctx.Cli.Bool("no-build") {
		BuildBinary(ctx)
	}

	binaryPath := typienv.Binary(ctx.Typical.BinaryNameOrDefault())

	log.Infof("Run the Binary '%s'", binaryPath)
	return bash.Run(binaryPath, []string(ctx.Cli.Args())...)
}
