package app

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/labstack/echo"
	"github.com/typical-go/typical-rest-server/app/controller"
)

// Config server configuration
type Config struct {
	Address string `envconfig:"ADDRESS" required:"true"`
}

// Server server application
type Server struct {
	*echo.Echo
	Config
	bookController controller.BookController
}

// NewServer return instance of server
func NewServer(
	config Config,
	bookController controller.BookController,
) *Server {

	s := &Server{
		Echo:           echo.New(),
		Config:         config,
		bookController: bookController,
	}
	initMiddlewares(s)
	initRoutes(s)

	return s
}

// CRUDController CRUD Controller
func (s *Server) CRUDController(entity string, crud controller.CRUDController) {
	s.GET(fmt.Sprintf("/%s", entity), crud.List)
	s.POST(fmt.Sprintf("/%s", entity), crud.Create)
	s.GET(fmt.Sprintf("/%s/:id", entity), crud.Get)
	s.PUT(fmt.Sprintf("/%s", entity), crud.Update)
	s.DELETE(fmt.Sprintf("/%s/:id", entity), crud.Delete)
}

// Serve start serve http request
func (s *Server) Serve() error {
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

	return s.Start(s.Address)
}
