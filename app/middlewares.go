package app

import (
	"github.com/labstack/echo/middleware"
)

func initMiddlewares(s *Server) {
	s.Use(middleware.Logger())
	s.Use(middleware.Recover())
	// check list of middleware at https://echo.labstack.com/middleware
}

// Put custom middleware belows
// Example: https://echo.labstack.com/cookbook/middleware
