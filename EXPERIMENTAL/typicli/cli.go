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
			if destroyer, ok := c.Obj.(typiobj.Destroyer); ok {
				if err = typiobj.Destroy(di, destroyer); err != nil {
					fmt.Println("Error: " + err.Error())
					os.Exit(1)
				}
				os.Exit(0)
			}
		}()
		if provider, ok := c.Obj.(typiobj.Provider); ok {
			if err = typiobj.Provide(di, provider); err != nil {
				return
			}
		}
		if preparer, ok := c.Obj.(typiobj.Preparer); ok {
			if err = typiobj.Prepare(di, preparer); err != nil {
				return
			}
		}
		return di.Invoke(fn)
	}
}
