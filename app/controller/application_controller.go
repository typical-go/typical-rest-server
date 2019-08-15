package controller

import (
	"database/sql"

	"github.com/typical-go/typical-rest-server/pkg/echokit"

	"github.com/labstack/echo"
)

// ApplicationController to handle API related with application itself
type ApplicationController struct {
	conn *sql.DB
}

// NewApplicationController return new instance of ApplicationController
func NewApplicationController(conn *sql.DB) *ApplicationController {
	return &ApplicationController{conn: conn}
}

// Route to define API Route
func (c *ApplicationController) Route(e *echo.Echo) {
	e.Any("application/health", c.Health)
}

// Health end point for health check
func (c *ApplicationController) Health(ctx echo.Context) error {
	return echokit.NewHealthCheck().
		Add("database", c.conn.Ping()).
		Send(ctx)
}
