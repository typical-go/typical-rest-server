package buildtool

import (
	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/typicli"
	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/typictx"
	"github.com/urfave/cli"
)

// Commands return list of command
func Commands(ctx *typictx.Context) (cmds []cli.Command) {
	for _, module := range ctx.AllModule() {
		if cmd := command(ctx, module); cmd != nil {
			cmds = append(cmds, *cmd)
		}
	}
	return
}

func command(ctx *typictx.Context, module interface{}) *cli.Command {
	if commander, ok := module.(typicli.BuildCommander); ok {
		cmd := commander.BuildCommand(&typicli.ContextCli{
			Context: ctx,
		})
		return &cmd
	}
	if commander, ok := module.(typicli.Commander); ok {
		cmd := commander.Command(&typicli.Cli{
			Obj: module,
		})
		return &cmd
	}
	return nil
}
