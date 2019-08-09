package typictx

import (
	"time"

	"gopkg.in/urfave/cli.v1"
)

// ActionFunc represented the action
type ActionFunc func(*ActionContext) error

// Action for cli
type Action interface {
	Start(*ActionContext) error
}

// ActionContext contain typical context and cli context
type ActionContext struct {
	Typical Context
	Cli     *cli.Context
}

// Deadline implementation
func (*ActionContext) Deadline() (deadline time.Time, ok bool) {
	return
}

// Done implementation
func (*ActionContext) Done() <-chan struct{} {
	return nil
}

// Err implementation
func (*ActionContext) Err() error {
	return nil
}

// Value implementation
func (*ActionContext) Value(key interface{}) interface{} {
	return nil
}
