package app

import (
	"context"
	"time"

	"github.com/labstack/echo"
	"github.com/typical-go/typical-rest-server/app/controller"
	"github.com/typical-go/typical-rest-server/config"
)

// Server server application
type Server struct {
	*echo.Echo
	*config.AppConfig
	bookController        controller.BookController
	applicationController controller.ApplicationController
}

// NewServer return instance of server
func NewServer(
	config *config.AppConfig,
	bookController controller.BookController,
	applicationController controller.ApplicationController,
) *Server {

	s := &Server{
		Echo:                  echo.New(),
		AppConfig:             config,
		bookController:        bookController,
		applicationController: applicationController,
	}

	initLogger(s)
	initMiddlewares(s)
	initRoutes(s)

	return s
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
