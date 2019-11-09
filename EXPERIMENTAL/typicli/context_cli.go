package typicli

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/typictx"

	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/typiobj"
	"github.com/urfave/cli"
	"go.uber.org/dig"
)

// ContextCli implementation of CLI
type ContextCli struct {
	*typictx.Context
}

// Action to return action function
func (c *ContextCli) Action(fn interface{}) func(ctx *cli.Context) error {
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
			fmt.Print("\n\n\n[[Application stop]]\n")
			if err = c.shutdown(di); err != nil {
				fmt.Println("Error: " + err.Error())
				os.Exit(1)
			}
			os.Exit(0)
		}()
		if err = c.beforeStart(di); err != nil {
			return
		}
		return di.Invoke(fn)
	}
}

func (c *ContextCli) beforeStart(di *dig.Container) (err error) {
	if err = provide(di, c.Constructors...); err != nil {
		return
	}
	for _, module := range c.AllModule() {
		if provider, ok := module.(typiobj.Provider); ok {
			if err = provide(di, provider.Provide()...); err != nil {
				return
			}
		}
	}
	for _, module := range c.AllModule() {
		if preparer, ok := module.(typiobj.Preparer); ok {
			if err = invoke(di, preparer.Prepare()...); err != nil {
				return
			}
		}
	}
	return
}

func (c *ContextCli) shutdown(di *dig.Container) (err error) {
	for _, module := range c.AllModule() {
		if destroyer, ok := module.(typiobj.Destroyer); ok {
			if err = invoke(di, destroyer.Destroy()...); err != nil {
				return
			}
		}
	}
	return
}

func invoke(di *dig.Container, fns ...interface{}) (err error) {
	for _, fn := range fns {
		if err = di.Invoke(fn); err != nil {
			return
		}
	}
	return
}

func provide(di *dig.Container, fns ...interface{}) (err error) {
	for _, fn := range fns {
		if err = di.Provide(fn); err != nil {
			return
		}
	}
	return
}
