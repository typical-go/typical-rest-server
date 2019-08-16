package typictx

import (
	"github.com/urfave/cli"
)

// ActionFunc represented the action
type ActionFunc func(*ActionContext) error

// Action for cli
type Action interface {
	Start(*ActionContext) error
}

// ActionContext contain typical context and cli context
type ActionContext struct {
	*Context
	Cli *cli.Context
}

// ActionFunction to convert function to ActionFunction
func ActionFunction(function interface{}) ActionFunc {
	return func(ctx *ActionContext) (err error) {
		return ctx.Container().Invoke(function)
	}
}
