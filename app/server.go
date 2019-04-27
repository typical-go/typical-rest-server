package app

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/typical-go/typical-rest-server/app/cntrl"
	"github.com/typical-go/typical-rest-server/config"
	"github.com/labstack/echo"
)

type server struct {
	*echo.Echo
	address        string
	bookController cntrl.BookController
}

func newServer(
	conf config.Config,
	bookController cntrl.BookController,
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

func (s *server) CRUDController(entity string, crud cntrl.CRUDController) {
	s.GET(fmt.Sprintf("/%s", entity), crud.List)
	s.POST(fmt.Sprintf("/%s", entity), crud.Create)
	s.GET(fmt.Sprintf("/%s/:id", entity), crud.Get)
	s.PUT(fmt.Sprintf("/%s", entity), crud.Update)
	s.DELETE(fmt.Sprintf("/%s/:id", entity), crud.Delete)
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
