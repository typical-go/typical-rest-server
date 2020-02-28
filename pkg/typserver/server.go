package typserver

import (
	"context"
	"fmt"
	"time"

	logrusmiddleware "github.com/bakatz/echo-logrusmiddleware"
	"github.com/labstack/echo"
	log "github.com/sirupsen/logrus"
)

// Server of rest
type Server struct {
	*echo.Echo
	logMiddleware echo.MiddlewareFunc
}

// New server instance
func New() *Server {
	e := echo.New()
	e.HideBanner = true
	e.Logger = logrusmiddleware.Logger{Logger: log.StandardLogger()}

	return &Server{
		Echo: e,
	}
}

// Shutdown the server
func Shutdown(s *Server) error {
	fmt.Println("Server shutting down")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	return s.Echo.Shutdown(ctx)
}

// SetDebug to set debug
func (s *Server) SetDebug(debug bool) {
	s.Debug = debug
	if debug {
		log.SetLevel(log.DebugLevel)
		log.SetFormatter(&log.TextFormatter{})
		s.logMiddleware = logrusmiddleware.HookWithConfig(logrusmiddleware.Config{
			IncludeRequestBodies:  true,
			IncludeResponseBodies: true,
		})
	} else {
		log.SetLevel(log.WarnLevel)
		log.SetFormatter(&log.JSONFormatter{})
		s.logMiddleware = logrusmiddleware.HookWithConfig(logrusmiddleware.Config{})
	}
}

// Start the server
func (s *Server) Start(addr string) error {
	s.Echo.Use(s.logMiddleware)
	return s.Echo.Start(addr)
}
