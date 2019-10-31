package typiobj

import (
	"github.com/urfave/cli"
	"go.uber.org/dig"
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

// CliAction to return cli action
func CliAction(p interface{}, fn interface{}) func(ctx *cli.Context) error {
	return func(ctx *cli.Context) (err error) {
		c := dig.New()
		defer func() {
			if destroyer, ok := p.(Destroyer); ok {
				if err = Destroy(c, destroyer); err != nil {
					return
				}
			}
		}()
		if provider, ok := p.(Provider); ok {
			if err = Provide(c, provider); err != nil {
				return
			}
		}
		return c.Invoke(fn)
	}
}
