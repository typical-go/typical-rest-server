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
	defer a.Close(di)
	gracefulStop := make(chan os.Signal)
	signal.Notify(gracefulStop, syscall.SIGTERM)
	signal.Notify(gracefulStop, syscall.SIGINT)
	if err = typiobj.Provide(di, a); err != nil {
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
		a.Close(di)
	}()
	runner := a.Application.(typiobj.Runner)
	return runner.Run(di)
}

func (a application) Provide() (constructors []interface{}) {
	constructors = append(constructors, a.Constructors...)
	constructors = append(constructors, a.Modules.Provide()...)
	if provider, ok := a.Application.(typiobj.Provider); ok {
		constructors = append(constructors, provider.Provide()...)
	}
	return
}

func (a application) Destroy() (destructors []interface{}) {
	destructors = append(destructors, a.Modules.Destroy()...)
	if destroyer, ok := a.Application.(typiobj.Destroyer); ok {
		destructors = append(destructors, destroyer.Destroy()...)
	}
	return
}

func (a application) Close(c *dig.Container) (err error) {
	return typiobj.Destroy(c, a)
}
