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
func (c ContextCli) Action(fn interface{}) func(ctx *cli.Context) error {
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
			if err = typiobj.Destroy(di, c); err != nil {
				fmt.Println("Error: " + err.Error())
				os.Exit(1)
			}
			os.Exit(0)
		}()
		if err = typiobj.Provide(di, c); err != nil {
			return
		}
		if err = typiobj.Prepare(di, c); err != nil {
			return
		}
		return di.Invoke(fn)
	}
}
