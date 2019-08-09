package app

import (
	"github.com/labstack/echo/middleware"
)

func initMiddlewares(s *Server) {
	s.Use(middleware.Recover())

}
