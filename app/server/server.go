package server

import "github.com/labstack/echo"

// Server server
type Server struct {
	echo.Echo
}

// New create new server
func New() *Server {
	server := &Server{
		Echo: *echo.New(),
	}

	initMiddlewares(server)
	initRoutes(server)

	return server
}
