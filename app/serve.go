package app

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/imantung/typical-go-server/config"
)

func serve(s *server, conf config.Config) error {
	gracefulStop := make(chan os.Signal)
	signal.Notify(gracefulStop, syscall.SIGTERM)
	signal.Notify(gracefulStop, syscall.SIGINT)

	// gracefull shutdown
	go func() {
		<-gracefulStop
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		s.Shutdown(ctx)
	}()

	return s.Start(conf.Address)
}
