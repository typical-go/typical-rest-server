package controller

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo"
)

const (
	underConstrucionStatus = http.StatusServiceUnavailable
	invalidMessageStatus   = http.StatusBadRequest
	invalidIDStatus        = http.StatusBadRequest
	insertSuccessStatus    = http.StatusCreated
)

func underContruction(ctx echo.Context) error {
	return ctx.JSON(underConstrucionStatus, map[string]string{"message": "Under Construction"})
}

func invalidMessage(ctx echo.Context, err error) error {
	res := map[string]interface{}{}
	res["message"] = "Invalid Message"

	return ctx.JSON(invalidMessageStatus, res)
}

func invalidID(ctx echo.Context, err error) error {
	res := map[string]interface{}{}
	res["message"] = "Invalid ID"

	return ctx.JSON(invalidIDStatus, res)
}

func insertSuccess(ctx echo.Context, lastInsertID int64) error {
	res := map[string]interface{}{}
	res["message"] = fmt.Sprintf("Success insert new record #%d", lastInsertID)
	return ctx.JSON(insertSuccessStatus, res)
}
