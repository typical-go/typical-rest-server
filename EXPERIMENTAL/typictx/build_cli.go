package typictx

import (
	"github.com/urfave/cli"
)

// BuildCLI responsible to give command
type BuildCLI interface {
	BuildCommand(ctx *Context) cli.Command
}

// IsBuildCLI return true if object implementation of BuildCLI
func IsBuildCLI(obj interface{}) (ok bool) {
	_, ok = obj.(BuildCLI)
	return
}
