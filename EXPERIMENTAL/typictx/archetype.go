package typictx

import "gopkg.in/urfave/cli.v1"

type ArcheType interface {
	GetMockTargets() []string
	GetTestTargets() []string
	GetConstructors() []interface{}
	GetConfig() interface{}
	GetConfigPrefix() string
	GetCommands() []Command
	StartApplication(ctx StartContext)
}

type StartContext struct {
	Context
	CliContext *cli.Context
}
