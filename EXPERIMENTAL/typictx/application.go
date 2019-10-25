package typictx

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	log "github.com/sirupsen/logrus"
	"github.com/typical-go/typical-rest-server/pkg/utility/errkit"
)

// Application is represent the application
type Application struct {
	Config
	StartFunc   interface{}
	StopFunc    interface{}
	Commands    []Command
	Initiations []interface{}
}

// Start the action
func (a Application) Start(ctx *ActionContext) (err error) {
	log.Info("------------- Application Start -------------")
	defer log.Info("-------------- Application End --------------")
	gracefulStop := make(chan os.Signal)
	signal.Notify(gracefulStop, syscall.SIGTERM)
	signal.Notify(gracefulStop, syscall.SIGINT)
	for _, initiation := range a.Initiations {
		err = ctx.Invoke(initiation)
		if err != nil {
			return
		}
	}
	go func() { // gracefull shutdown
		<-gracefulStop
		fmt.Println() // NOTE: intentionally print new line after "^C"
		fmt.Println("Graceful Shutdown...")
		var errs errkit.Errors
		if a.StopFunc != nil {
			errs.Add(ctx.Invoke(a.StopFunc))
		}
		for _, module := range ctx.Modules {
			if module.CloseFunc != nil {
				errs.Add(ctx.Invoke(module.CloseFunc))
			}
		}
		err = errs
	}()
	if a.StartFunc != nil {
		err = ctx.Invoke(a.StartFunc)
	}
	return
}
