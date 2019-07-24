package typimain

import (
	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/typictx"
	"gopkg.in/urfave/cli.v1"
)

func runActionFunc(context typictx.Context, actionFunc typictx.ActionFunc) interface{} {
	return func(ctx *cli.Context) error {
		return actionFunc(typictx.ActionContext{
			Cli:     ctx,
			Typical: context,
		})

	}
}

func runAction(context typictx.Context, action typictx.Action) interface{} {
	return func(ctx *cli.Context) error {
		return action.Start(typictx.ActionContext{
			Cli:     ctx,
			Typical: context,
		})
	}
}
