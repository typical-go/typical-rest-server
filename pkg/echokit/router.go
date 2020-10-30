package echokit

import (
	"github.com/labstack/echo/v4"
)

type (
	// Server interface for echo.Echo and echo.Group
	Server interface {
		CONNECT(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
		DELETE(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
		GET(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
		HEAD(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
		OPTIONS(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
		PATCH(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
		POST(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
		PUT(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
		TRACE(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
		Any(path string, handler echo.HandlerFunc, middleware ...echo.MiddlewareFunc) []*echo.Route
		Add(method, path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
		Match(methods []string, path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) []*echo.Route
		Group(prefix string, m ...echo.MiddlewareFunc) *echo.Group
		Use(m ...echo.MiddlewareFunc)
	}
	// Router responsible to route
	Router interface {
		SetRoute(Server)
	}
	// SetRouteFn function SetRoute
	SetRouteFn func(Server)
	routerImpl struct {
		fn SetRouteFn
	}
)

var _ Server = (*echo.Echo)(nil)
var _ Server = (*echo.Group)(nil)

// SetRoute to server
func SetRoute(server Server, routers ...Router) {
	for _, router := range routers {
		router.SetRoute(server)
	}
}

//
// routerImpl
//

// NewRouter return new instance of Router
func NewRouter(fn SetRouteFn) Router {
	return &routerImpl{
		fn: fn,
	}
}

// SetRoute set route
func (r *routerImpl) SetRoute(server Server) {
	r.fn(server)
}
