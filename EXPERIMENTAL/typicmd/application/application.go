package application

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

type application struct {
	*typictx.Context
}

func (a application) Run(ctx *cli.Context) (err error) {
	di := dig.New()
	defer a.Destruct(di)
	gracefulStop := make(chan os.Signal)
	signal.Notify(gracefulStop, syscall.SIGTERM)
	signal.Notify(gracefulStop, syscall.SIGINT)
	if err = a.Construct(di); err != nil {
		return
	}
	// TODO: create prepare function
	// for _, initiation := range a.Initiations {
	// 	if err = di.Invoke(initiation); err != nil {
	// 		return
	// 	}
	// }
	go func() {
		<-gracefulStop
		fmt.Println("\n\n\nGraceful Shutdown...")
		a.Destruct(di)
	}()
	runner := a.Application.(typiobj.Runner)
	return runner.Run(di)
}
