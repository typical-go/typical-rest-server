package typiobj

import (
	"github.com/urfave/cli"
	"go.uber.org/dig"
)

// Action to return cli action
func Action(p interface{}, fn interface{}) func(ctx *cli.Context) error {
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
