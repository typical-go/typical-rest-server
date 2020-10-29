package app

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/go-redis/redis"
	"github.com/labstack/echo/v4"
	"github.com/typical-go/typical-go/pkg/typgo"
	"github.com/typical-go/typical-rest-server/pkg/typrest"
	"go.uber.org/dig"
)

type (
	// HealthCheck ...
	HealthCheck struct {
		dig.In
		PG    *sql.DB `name:"pg"`
		MySQL *sql.DB `name:"mysql"`
		Redis *redis.Client
	}
)

// Handle echo function
func (h *HealthCheck) Handle(ec echo.Context) error {
	health := typrest.HealthMap{
		"postgres": h.PG.Ping(),
		"mysql":    h.MySQL.Ping(),
		"redis":    h.Redis.Ping().Err(),
	}

	status, ok := health.Status()
	return ec.JSON(h.httpStatus(ok), h.response(status))
}

func (h *HealthCheck) httpStatus(ok bool) int {
	httpStatus := http.StatusServiceUnavailable
	if ok {
		httpStatus = http.StatusOK
	}
	return httpStatus
}

func (h *HealthCheck) response(status map[string]string) map[string]interface{} {
	return map[string]interface{}{
		"name":   fmt.Sprintf("%s (%s)", typgo.ProjectName, typgo.ProjectVersion),
		"status": status,
	}
}
