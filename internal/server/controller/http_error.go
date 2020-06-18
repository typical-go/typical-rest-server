package controller

import (
	"database/sql"
	"net/http"

	"github.com/typical-go/typical-rest-server/pkg/errvalid"

	"github.com/labstack/echo"
)

func httpError(err error) *echo.HTTPError {
	if err == sql.ErrNoRows {
		return echo.NewHTTPError(http.StatusNotFound)
	}

	if errvalid.Check(err) {
		return echo.NewHTTPError(
			http.StatusUnprocessableEntity,
			errvalid.Message(err),
		)
	}

	return echo.NewHTTPError(
		http.StatusInternalServerError,
		err.Error(),
	)
}
