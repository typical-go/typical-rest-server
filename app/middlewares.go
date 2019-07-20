package app

import (
	logrusmiddleware "github.com/bakatz/echo-logrusmiddleware"
	"github.com/labstack/echo/middleware"
)

func initMiddlewares(s *Server) {

	s.Use(middleware.Recover())
	s.Use(logrusmiddleware.HookWithConfig(logrusmiddleware.Config{
		IncludeRequestBodies:  true,
		IncludeResponseBodies: true,
	}))

}
