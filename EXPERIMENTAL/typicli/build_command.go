package typicli

import (
	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/typictx"
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

// BuildCommands return list of build command
func BuildCommands(ctx *typictx.Context) (cmds []cli.Command) {
	for _, module := range ctx.AllModule() {
		if commander, ok := module.(BuildCommander); ok {
			cmds = append(cmds, commander.BuildCommand(&ContextCli{
				Context: ctx,
			}))
		}
	}
	return
}
