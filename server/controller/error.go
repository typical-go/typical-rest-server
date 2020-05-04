package controller

import (
	"database/sql"
	"net/http"

	"github.com/labstack/echo"
)

func httpError(err error) *echo.HTTPError {
	switch err {
	case sql.ErrNoRows:
		return echo.NewHTTPError(http.StatusNotFound)
	default:
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
}
