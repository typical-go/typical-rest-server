package typiobj

import (
	"github.com/urfave/cli"
	"go.uber.org/dig"
)

// CliAction to return cli action
func CliAction(p interface{}, fn interface{}) func(ctx *cli.Context) error {
	return func(ctx *cli.Context) (err error) {
		c := dig.New()
		defer func() {
			if destructor, ok := p.(Destructor); ok {
				destructor.Destruct(c)
			}
		}()
		if provider, ok := p.(Provider); ok {
			for _, constructor := range provider.Provide() {
				if err = c.Provide(constructor); err != nil {
					return
				}
			}
		}
		return c.Invoke(fn)
	}
}
