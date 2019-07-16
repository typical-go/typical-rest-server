package typictx

import (
	"os"
	"os/signal"
	"syscall"
)

type RunAction struct {
	StartFunc interface{}
	StopFunc  interface{}
}

// Start the action
func (a RunAction) Start(ctx ActionContext) (err error) {
	container := ctx.Container()

	if a.StopFunc != nil {
		gracefulStop := make(chan os.Signal)
		signal.Notify(gracefulStop, syscall.SIGTERM)
		signal.Notify(gracefulStop, syscall.SIGINT)

		// gracefull shutdown
		go func() {
			<-gracefulStop
			err = container.Invoke(a.StopFunc)
		}()
	}

	if a.StartFunc != nil {
		err = container.Invoke(a.StartFunc)
	}

	return
}
