package typictx

import "github.com/urfave/cli"

// AppCLI return command
type AppCLI interface {
	AppCommands(ctx *Context) []cli.Command
}
