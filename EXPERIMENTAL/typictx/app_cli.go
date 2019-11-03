package typictx

import "github.com/urfave/cli"

// AppCLI return command
type AppCLI interface {
	AppCommands(ctx *Context) []cli.Command
}

// IsAppCLI return true if object implementation of AppCLI
func IsAppCLI(obj interface{}) (ok bool) {
	_, ok = obj.(AppCLI)
	return
}
