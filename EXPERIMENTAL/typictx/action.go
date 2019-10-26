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
