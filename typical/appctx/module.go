package appctx

import (
	"gopkg.in/urfave/cli.v1"
)

// Module in typical-go applicaiton
type Module interface {
	Config() interface{}
	ConfigPrefix() string
	Constructors() []interface{}
	Command() cli.Command
}
