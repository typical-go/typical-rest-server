package typictx

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	log "github.com/sirupsen/logrus"
)

// Application is represent the application
type Application struct {
	Config      Config
	StartFunc   interface{}
	StopFunc    interface{}
	Commands    []*Command
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
		if err = ctx.Invoke(initiation); err != nil {
			return
		}
	}
	go func() { // gracefull shutdown
		<-gracefulStop
		fmt.Println() // NOTE: intentionally print new line after "^C"
		fmt.Println("Graceful Shutdown...")
		err = ctx.Destruct(ctx.Container)
	}()
	if a.StartFunc != nil {
		err = ctx.Invoke(a.StartFunc)
	}
	return
}

// Configure return configuration
func (a Application) Configure() Config {
	return a.Config
}
