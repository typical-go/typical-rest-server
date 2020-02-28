package typserver

import (
	"context"
	"fmt"
	"net/http"
	"time"

	logrusmiddleware "github.com/bakatz/echo-logrusmiddleware"
	"github.com/labstack/echo"
	log "github.com/sirupsen/logrus"
)

// Server of rest
type Server struct {
	*echo.Echo
	logMiddleware       echo.MiddlewareFunc
	healthChecker       map[string]func() error
	healthCheckEndpoint string
}

// New server instance
func New() *Server {
	e := echo.New()
	e.HideBanner = true
	e.Logger = logrusmiddleware.Logger{Logger: log.StandardLogger()}

	return &Server{
		Echo:                e,
		healthChecker:       make(map[string]func() error),
		healthCheckEndpoint: "application/health",
	}
}

// Shutdown the server
func Shutdown(s *Server) error {
	fmt.Println("Server shutting down")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	return s.Echo.Shutdown(ctx)
}

// WithHealthCheckEndpoint return Server with new health check endpoint
func (s *Server) WithHealthCheckEndpoint(endpoint string) *Server {
	s.healthCheckEndpoint = endpoint
	return s
}

// PutHealthChecker to put health check function
func (s *Server) PutHealthChecker(name string, fn func() error) {
	s.healthChecker[name] = fn
}

// Register controller
func (s *Server) Register(cntrl Controller) {
	cntrl.SetRoute(s.Echo)
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
	s.Echo.Any(s.healthCheckEndpoint, s.healthCheckHandler)
	return s.Echo.Start(addr)
}

func (s *Server) healthCheckHandler(ctx echo.Context) error {
	status := http.StatusOK
	message := make(map[string]string)

	for name, fn := range s.healthChecker {
		if err := fn(); err != nil {
			message[name] = err.Error()
			status = http.StatusServiceUnavailable
		} else {
			message[name] = "OK"
		}
	}

	return ctx.JSON(status, message)
}
