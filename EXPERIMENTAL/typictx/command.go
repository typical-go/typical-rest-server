package typictx

import "github.com/urfave/cli"

// CommandLiner responsible to give command
type CommandLiner interface {
	CommandLine() cli.Command
}
