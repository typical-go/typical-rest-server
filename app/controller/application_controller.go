package controller

import (
	"database/sql"
	"net/http"

	"github.com/labstack/echo"
	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/typical-arc-rest/utility"
)

// ApplicationController to handle API related with application itself
type ApplicationController struct {
	conn *sql.DB
}

// NewApplicationController return new instance of ApplicationController
func NewApplicationController(conn *sql.DB) ApplicationController {
	return ApplicationController{
		conn: conn,
	}
}

// Health end point for health check
func (c *ApplicationController) Health(ctx echo.Context) error {
	healthCheck := utility.NewHealthCheck()
	healthCheck.Add("database", c.conn.Ping())

	var status int
	var message interface{}

	if ctx.QueryParam("full") != "" {
		message = healthCheck
	}

	if healthCheck.NotOK() {
		status = http.StatusServiceUnavailable
	} else {
		status = http.StatusOK
	}

	return ctx.JSON(status, message)

}
