package typicli

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/typiobj"
	"github.com/urfave/cli"
	"go.uber.org/dig"
)

// Cli for command line interface
type Cli struct {
	Obj interface{}
}

// Action to return action function
func (c Cli) Action(fn interface{}) func(ctx *cli.Context) error {
	return func(ctx *cli.Context) (err error) {
		di := dig.New()
		gracefulStop := make(chan os.Signal)
		signal.Notify(gracefulStop, syscall.SIGTERM)
		signal.Notify(gracefulStop, syscall.SIGINT)
		defer func() {
			gracefulStop <- syscall.SIGTERM
		}()
		go func() {
			<-gracefulStop
			if err = c.shutdown(di); err != nil {
				fmt.Println("Error: " + err.Error())
				os.Exit(1)
				return
			}
			os.Exit(0)
		}()
		if err = c.provideDependency(di); err != nil {
			return
		}
		if err = c.prepare(di); err != nil {
			return
		}
		return di.Invoke(fn)
	}
}

func (c Cli) provideDependency(di *dig.Container) (err error) {
	if provider, ok := c.Obj.(typiobj.Provider); ok {
		if err = provide(di, provider.Provide()...); err != nil {
			return
		}
	}
	return
}

func (c Cli) prepare(di *dig.Container) (err error) {
	if preparer, ok := c.Obj.(typiobj.Preparer); ok {
		if err = invoke(di, preparer.Prepare()...); err != nil {
			return
		}
	}
	return
}

func (c Cli) shutdown(di *dig.Container) (err error) {
	if destroyer, ok := c.Obj.(typiobj.Destroyer); ok {
		if err = invoke(di, destroyer.Destroy()...); err != nil {
			return
		}
	}
	return
}
