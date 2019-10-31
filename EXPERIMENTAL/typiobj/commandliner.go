package typiobj

import (
	"github.com/urfave/cli"
)

// CommandLiner responsible to give command
type CommandLiner interface {
	CommandLine() cli.Command
}

// IsCommandLiner return true if object implementation of CommandLiner
func IsCommandLiner(obj interface{}) (ok bool) {
	_, ok = obj.(CommandLiner)
	return
}
