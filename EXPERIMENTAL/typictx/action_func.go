package typictx

import "github.com/urfave/cli"

// ActionFunc represented the action
type ActionFunc func(*ActionContext) error

// ActionFunction to convert function to ActionFunction
func ActionFunction(function interface{}) ActionFunc {
	return func(ctx *ActionContext) (err error) {
		return ctx.Invoke(function)
	}
}

// CommandFunction to convert to command function
func (f ActionFunc) CommandFunction(context *Context) interface{} {
	return func(ctx *cli.Context) error {
		return f(&ActionContext{
			Cli:     ctx,
			Context: context,
		})
	}
}
