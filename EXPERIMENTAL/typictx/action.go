package typictx

import "gopkg.in/urfave/cli.v1"

type Action interface {
	Start(ActionContext) error
}

type ActionContext struct {
	Context
	CliContext *cli.Context
}
