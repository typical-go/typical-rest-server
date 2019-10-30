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
		if constructor, ok := p.(Constructor); ok {
			if err = constructor.Construct(c); err != nil {
				return
			}
		}
		return c.Invoke(fn)
	}
}
