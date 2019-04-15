package app

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/imantung/typical-go-server/app/controller"
	"github.com/imantung/typical-go-server/config"
	"github.com/labstack/echo"
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

type server struct {
	*echo.Echo
	bookController controller.BookController
}

func newServer(bookController controller.BookController) *server {
	s := &server{
		Echo:           echo.New(),
		bookController: bookController,
	}

	initMiddlewares(s)
	initRoutes(s)

	return s
}

func (s *server) CRUD(entity string, crud controller.CRUD) {
	crud.RegisterTo(entity, s.Echo)
}
