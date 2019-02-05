package server

import "github.com/labstack/echo/middleware"

func initMiddlewares(server *Server) {
	server.Use(middleware.Logger())
	server.Use(middleware.Recover())
	// check list of middleware at https://echo.labstack.com/middleware
}

// Put custom middleware belows
// Example: https://echo.labstack.com/cookbook/middleware
