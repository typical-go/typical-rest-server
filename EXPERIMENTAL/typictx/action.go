package typictx

import "gopkg.in/urfave/cli.v1"

// ActionFunc represented the action
type ActionFunc func(ActionContext) error

// Action for cli
type Action interface {
	Start(ActionContext) error
}

// ActionContext contain typical context and cli context
type ActionContext struct {
	Context
	CliContext *cli.Context
}
