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
	"github.com/typical-go/typical-rest-server/pkg/restkit"
	"go.uber.org/dig"
)

type (
	// HealthCheck ...
	HealthCheck struct {
		dig.In
		PG    *sql.DB `name:"pg"`
		Cache *cachekit.Store
	}
)

//
// HealthCheck
//

// Handle echo function
func (h *HealthCheck) Handle(ec echo.Context) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	health := restkit.HealthMap{
		"postgres": h.PG.Ping(),
		"cache":    h.Cache.Ping(ctx).Err(),
	}

	// See: http://www.w3.org/Protocols/rfc2616/rfc2616-sec14.html#sec14.21
	ec.Response().Header().Set("Expires", "0")

	// See: http://www.w3.org/Protocols/rfc2616/rfc2616-sec14.html#sec14.32
	ec.Response().Header().Set("Pragma", "no-cache")

	// See: http://www.w3.org/Protocols/rfc2616/rfc2616-sec14.html#sec14.9
	ec.Response().Header().Set(
		"Cache-Control",
		"no-cache, no-store, must-revalidate",
	)

	status, ok := health.Status()
	code := h.httpStatus(ok)

	if ec.Request().Method == http.MethodHead {
		return ec.NoContent(code)
	} else {
		return ec.JSON(code, h.response(status))
	}
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
