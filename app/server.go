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

type server struct {
	*echo.Echo
	address        string
	bookController controller.BookController
}

func newServer(
	conf config.Config,
	bookController controller.BookController,
) *server {

	s := &server{
		Echo:           echo.New(),
		address:        conf.Address,
		bookController: bookController,
	}

	initMiddlewares(s)
	initRoutes(s)

	return s
}

func (s *server) CRUD(entity string, crud controller.CRUD) {
	crud.RegisterTo(entity, s.Echo)
}

func (s *server) Serve() error {
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

	return s.Start(s.address)
}
