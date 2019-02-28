package controller

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/labstack/echo"
)

const (
	underConstrucionStatus = http.StatusServiceUnavailable
	invalidMessageStatus   = http.StatusBadRequest
	insertSuccessStatus    = http.StatusCreated
)

func underContruction(c echo.Context) error {
	return c.JSON(underConstrucionStatus, map[string]string{"message": "Under Construction"})
}

func invalidMessage(c echo.Context, err error) error {
	res := map[string]interface{}{}
	res["message"] = "Invalid Message"

	return c.JSON(invalidMessageStatus, res)
}

func insertSuccess(c echo.Context, result sql.Result) error {
	lastInsertID, _ := result.LastInsertId()
	res := map[string]interface{}{}
	res["message"] = fmt.Sprintf("Success insert new record #%d", lastInsertID)
	return c.JSON(insertSuccessStatus, res)
}
