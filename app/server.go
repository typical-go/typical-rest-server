package app

import (
	"github.com/labstack/echo"
	"github.com/typical-go/typical-rest-server/app/controller"
)

// Server server application
type Server struct {
	*echo.Echo
	bookController        controller.BookController
	applicationController controller.ApplicationController
}

// NewServer return instance of server
func NewServer(
	echoServer *echo.Echo,
	bookController controller.BookController,
	applicationController controller.ApplicationController,
) *Server {

	s := &Server{
		Echo:                  echoServer,
		bookController:        bookController,
		applicationController: applicationController,
	}

	initMiddlewares(s)
	initRoutes(s)

	return s
}
