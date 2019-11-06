package typicli

import "github.com/urfave/cli"

// Commander responsible to command
type Commander interface {
	Command(c *Cli) cli.Command
}

// IsCommander return true if obj implement commander
func IsCommander(obj interface{}) (ok bool) {
	_, ok = obj.(Commander)
	return
}
