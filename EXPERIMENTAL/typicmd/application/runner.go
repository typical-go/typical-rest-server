package application

import (
	"os"
	"os/signal"
	"syscall"

	log "github.com/sirupsen/logrus"
	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/typictx"
	"github.com/urfave/cli"
	"go.uber.org/dig"
)

type runner struct {
	*typictx.Context
	action interface{}
}

func (a runner) Run(ctx *cli.Context) (err error) {
	log.Info("------------- Application Start -------------")
	defer log.Info("-------------- Application End --------------")
	gracefulStop := make(chan os.Signal)
	signal.Notify(gracefulStop, syscall.SIGTERM)
	signal.Notify(gracefulStop, syscall.SIGINT)
	di := dig.New()
	if err = a.Construct(di); err != nil {
		return
	}
	// TODO: create prepare function
	for _, initiation := range a.Initiations {
		if err = di.Invoke(initiation); err != nil {
			return
		}
	}
	go func() { // gracefull shutdown
		<-gracefulStop
		log.Println("\nGraceful Shutdown...")
		err = a.Destruct(di)
	}()
	if err = di.Invoke(a.action); err != nil {
		return
	}
	return
}
