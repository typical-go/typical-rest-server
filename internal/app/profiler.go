package app

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/typical-go/typical-go/pkg/typgo"
	"github.com/typical-go/typical-rest-server/pkg/cachekit"
	"github.com/typical-go/typical-rest-server/pkg/typrest"
	"go.uber.org/dig"
)

type (
	// HealthCheck ...
	HealthCheck struct {
		dig.In
		PG    *sql.DB `name:"pg"`
		MySQL *sql.DB `name:"mysql"`
		Cache *cachekit.Store
	}
)

const (
	healthCheckPath = "/application/health"
)

// SetProfiler ...
func SetProfiler(e *echo.Echo, hc HealthCheck) {
	e.GET(healthCheckPath, hc.Handle)
	e.HEAD(healthCheckPath, hc.Handle)
	e.GET("/debug/*", echo.WrapHandler(http.DefaultServeMux))
	e.GET("/debug/*/*", echo.WrapHandler(http.DefaultServeMux))
}

//
// HealthCheck
//

// Handle echo function
func (h *HealthCheck) Handle(ec echo.Context) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	health := typrest.HealthMap{
		"postgres": h.PG.Ping(),
		"mysql":    h.MySQL.Ping(),
		"cache":    h.Cache.Ping(ctx).Err(),
	}

	status, ok := health.Status()
	return ec.JSON(h.httpStatus(ok), h.response(status))
}

func (h *HealthCheck) httpStatus(ok bool) int {
	if ok {
		return http.StatusOK
	}
	return http.StatusServiceUnavailable
}

func (h *HealthCheck) response(status map[string]string) map[string]interface{} {
	return map[string]interface{}{
		"name":   fmt.Sprintf("%s (%s)", typgo.ProjectName, typgo.ProjectVersion),
		"status": status,
	}
}
