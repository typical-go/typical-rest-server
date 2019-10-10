package responsekit

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo"
)

// Message is general response
type Message struct {
	Message string `json:"message"`
}

// InvalidRequest is response for invalid request
func InvalidRequest(ctx echo.Context, err error) error {
	return ctx.JSON(http.StatusBadRequest, Message{"Invalid Request"})
}

// InvalidID is response for invalid ID
func InvalidID(ctx echo.Context, err error) error {
	return ctx.JSON(http.StatusBadRequest, Message{"Invalid ID"})
}

// NotFound is response for not found
func NotFound(ctx echo.Context, id int64) error {
	msg := fmt.Sprintf("#%d not found", id)
	return ctx.JSON(http.StatusNotFound, Message{msg})
}

// InsertSuccess is response for insert success
func InsertSuccess(ctx echo.Context, lastInsertID int64) error {
	msg := fmt.Sprintf("Success insert new record #%d", lastInsertID)
	return ctx.JSON(http.StatusCreated, Message{msg})
}

// DeleteSuccess is response for delete success
func DeleteSuccess(ctx echo.Context, id int64) error {
	msg := fmt.Sprintf("Delete #%d done", id)
	return ctx.JSON(http.StatusOK, Message{msg})
}

// UpdateSuccess is response for update success
func UpdateSuccess(ctx echo.Context, id int64) error {
	msg := fmt.Sprintf("Update #%d success", id)
	return ctx.JSON(http.StatusOK, Message{msg})
}
