package app

import (
	"context"
	"fmt"
	"time"

	"github.com/labstack/echo"
	"github.com/typical-go/typical-rest-server/app/controller"
	"github.com/typical-go/typical-rest-server/config"
)

// Server server application
type Server struct {
	*echo.Echo
	config.AppConfig
	bookController controller.BookController
}

// NewServer return instance of server
func NewServer(
	config config.AppConfig,
	bookController controller.BookController,
) *Server {

	s := &Server{
		Echo:           echo.New(),
		AppConfig:      config,
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

// Serve to start to serve
func (s *Server) Serve() error {
	return s.Echo.Start(s.Address)
}

// Shutdown the server
func (s *Server) Shutdown() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	return s.Echo.Shutdown(ctx)
}
