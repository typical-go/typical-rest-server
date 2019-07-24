package typictx

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/internal/util"
)

// MainAction hadling main process action
type MainAction struct {
	StartFunc interface{}
	StopFunc  interface{}
}

// Start the action
func (a MainAction) Start(ctx ActionContext) (err error) {
	container := ctx.Typical.Container()

	if a.StopFunc != nil {
		gracefulStop := make(chan os.Signal)
		signal.Notify(gracefulStop, syscall.SIGTERM)
		signal.Notify(gracefulStop, syscall.SIGINT)

		// gracefull shutdown
		go func() {
			<-gracefulStop

			// NOTE: intentionally print new line after "^C"
			fmt.Println()

			var errs util.Errors
			errs.Add(container.Invoke(a.StopFunc))

			for _, module := range ctx.Typical.Modules {
				if module.CloseFunc != nil {
					errs.Add(container.Invoke(module.CloseFunc))
				}
			}

			err = errs
		}()
	}

	if a.StartFunc != nil {
		err = container.Invoke(a.StartFunc)
	}

	return
}
