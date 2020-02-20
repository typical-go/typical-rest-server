package controller

import (
	"github.com/labstack/echo"
	"github.com/typical-go/typical-rest-server/pkg/serverkit"
	"github.com/typical-go/typical-rest-server/pkg/typpostgres"
	"go.uber.org/dig"
)

// AppCntrl handle API related with application itself
type AppCntrl struct {
	dig.In
	Postgres *typpostgres.DB
}

// Route to define API Route
func (c *AppCntrl) Route(e *echo.Echo) {
	e.Any("application/health", c.Health)
}

// Health end point for health check
func (c *AppCntrl) Health(ctx echo.Context) error {
	return serverkit.NewHealthCheck().
		Add("postgres", c.Postgres.Ping()).
		Send(ctx)
}
