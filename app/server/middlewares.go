package server

import "github.com/labstack/echo/middleware"

func initMiddlewares(server *Server) {
	server.Use(middleware.Logger())
	server.Use(middleware.Recover())
}
