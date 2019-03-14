package controller

import (
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

func insertSuccess(c echo.Context, lastInsertID int64) error {
	res := map[string]interface{}{}
	res["message"] = fmt.Sprintf("Success insert new record #%d", lastInsertID)
	return c.JSON(insertSuccessStatus, res)
}
