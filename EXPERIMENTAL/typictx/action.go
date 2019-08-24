package typictx

import (
	"github.com/urfave/cli"
)

// Action for cli
type Action interface {
	Start(*ActionContext) error
}

// ActionContext contain typical context and cli context
type ActionContext struct {
	*Context
	Cli *cli.Context
}

// ActionCommandFunction to get command function fo action
func ActionCommandFunction(context *Context, action Action) interface{} {
	return func(ctx *cli.Context) error {
		return action.Start(&ActionContext{
			Cli:     ctx,
			Context: context,
		})
	}
}
