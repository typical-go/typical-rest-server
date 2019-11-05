package typicli

import "github.com/urfave/cli"

// AppCommander return command
type AppCommander interface {
	AppCommands(c *ContextCli) []cli.Command
}

// IsAppCommander return true if object implementation of AppCLI
func IsAppCommander(obj interface{}) (ok bool) {
	_, ok = obj.(AppCommander)
	return
}
