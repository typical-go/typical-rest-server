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
	Context
	Cli *cli.Context
}
