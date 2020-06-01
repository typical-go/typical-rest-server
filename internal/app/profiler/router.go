package profiler

import (
	"database/sql"

	"github.com/go-redis/redis"
	"github.com/typical-go/typical-rest-server/pkg/echokit"
	"go.uber.org/dig"
)

// Router for profiler
type Router struct {
	dig.In
	PG    *sql.DB
	Redis *redis.Client
}

var _ echokit.Router = (*Router)(nil)

// Route to profiler api
func (h *Router) Route(e echokit.Server) error {
	hc := echokit.HealthCheck{
		"postgres": h.PG.Ping,
		"redis":    h.Redis.Ping().Err,
	}

	e.Any("application/health", hc.JSON)
	return nil
}
