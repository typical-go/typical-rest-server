package echokit

import (
	"github.com/labstack/echo"
)

var (
	_ Server = (*echo.Echo)(nil)
	_ Server = (*echo.Group)(nil)
)

type (
	// Router responsible to route
	Router interface {
		Route(Server) error
	}
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
)

// SetRoute to server
func SetRoute(server Server, routers ...Router) (err error) {
	for _, router := range routers {
		if err = router.Route(server); err != nil {
			return
		}
	}
	return
}
