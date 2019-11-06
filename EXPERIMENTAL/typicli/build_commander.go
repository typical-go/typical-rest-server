package typicli

import (
	"github.com/urfave/cli"
)

// BuildCommander responsible to give command
type BuildCommander interface {
	BuildCommand(c *ContextCli) cli.Command
}

// IsBuildCommander return true if object implementation of BuildCLI
func IsBuildCommander(obj interface{}) (ok bool) {
	_, ok = obj.(BuildCommander)
	return
}
