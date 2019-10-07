package controller

import (
	"database/sql"

	"github.com/labstack/echo"
	"github.com/typical-go/typical-rest-server/pkg/utility/echokit"
	"go.uber.org/dig"
)

// AppCntrl handle API related with application itself
type AppCntrl struct {
	dig.In
	MYSQL *sql.DB
}

// Route to define API Route
func (c *AppCntrl) Route(e *echo.Echo) {
	e.Any("application/health", c.Health)
}

// Health end point for health check
func (c *AppCntrl) Health(ctx echo.Context) error {
	return echokit.NewHealthCheck().
		Add("database", c.MYSQL.Ping()).
		Send(ctx)
}
