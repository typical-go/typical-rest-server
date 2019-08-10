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
	container := ctx.Container()

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
			errs.Add(container.Invoke(a.StopFunc))
		}

		for _, module := range ctx.Modules {
			if module.CloseFunc != nil {
				errs.Add(container.Invoke(module.CloseFunc))
			}
		}

		err = errs
	}()

	if a.StartFunc != nil {
		err = container.Invoke(a.StartFunc)
	}

	return
}
