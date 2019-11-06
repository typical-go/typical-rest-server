package buildtool

import (
	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/typicli"
	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/typictx"
	"github.com/urfave/cli"
)

// Commands return list of command
// TODO: return command detail instead of list command for readme and metadata
func Commands(ctx *typictx.Context) (cmds []cli.Command) {
	for _, module := range ctx.AllModule() {
		cmds = append(cmds, command(ctx, module))
	}
	return
}

func command(ctx *typictx.Context, module interface{}) cli.Command {
	if commander, ok := module.(typicli.BuildCommander); ok {
		c := &typicli.ContextCli{Context: ctx}
		return commander.BuildCommand(c)
	}
	if commander, ok := module.(typicli.Commander); ok {
		c := &typicli.Cli{Obj: module}
		return commander.Command(c)
	}
	return cli.Command{}
}
