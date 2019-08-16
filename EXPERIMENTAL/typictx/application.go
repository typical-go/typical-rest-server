package typictx

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/typical-go/runn"
)

// Application is represent the application
type Application struct {
	StartFunc interface{}
	StopFunc  interface{}
	Commands  []Command
}

// Start the action
func (a Application) Start(ctx *ActionContext) (err error) {
	gracefulStop := make(chan os.Signal)
	signal.Notify(gracefulStop, syscall.SIGTERM)
	signal.Notify(gracefulStop, syscall.SIGINT)

	// gracefull shutdown
	go func() {
		<-gracefulStop

		// NOTE: intentionally print new line after "^C"
		fmt.Println()
		fmt.Println("Graceful Shutdown...")

		var errs runn.Errors
		if a.StopFunc != nil {
			errs.Add(ctx.Container().Invoke(a.StopFunc))
		}

		for _, module := range ctx.Modules {
			if module.CloseFunc != nil {
				errs.Add(ctx.Container().Invoke(module.CloseFunc))
			}
		}

		err = errs
	}()

	if a.StartFunc != nil {
		err = ctx.Container().Invoke(a.StartFunc)
	}

	return
}
