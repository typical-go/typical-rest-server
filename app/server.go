package app

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/imantung/typical-go-server/config"
	"github.com/labstack/echo"
)

type server struct {
	*echo.Echo
}

func newServer() *server {
	s := &server{
		Echo: echo.New(),
	}
	s.initMiddlewares()
	s.initRoutes()
	return s
}

func serve(s *server, conf config.Config) error {
	gracefulStop := make(chan os.Signal)
	signal.Notify(gracefulStop, syscall.SIGTERM)
	signal.Notify(gracefulStop, syscall.SIGINT)
	go func() { // gracefull shutdown
		<-gracefulStop
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		s.Shutdown(ctx)
	}()
	return s.Start(conf.Address)
}
