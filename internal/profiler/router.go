package profiler

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/typical-go/typical-rest-server/pkg/typrest"
	"go.uber.org/dig"

	// enable `/debug/vars`
	_ "expvar"

	// enable `/debug/pprof` API
	_ "net/http/pprof"
)

type (
	// Router for profiler
	Router struct {
		dig.In
		HealthCheck
	}
)

var _ typrest.Router = (*Router)(nil)

// SetRoute for profiler
func (r *Router) SetRoute(server typrest.Server) error {
	r.HealthCheck.SetRoute(server)
	server.GET("/debug/vars", echo.WrapHandler(http.DefaultServeMux))
	server.GET("/debug/pprof/*", echo.WrapHandler(http.DefaultServeMux))
	return nil
}
