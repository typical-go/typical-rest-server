package controller

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo"
	"go.uber.org/dig"
)

// BookControllerResponse is response of BookController
type BookControllerResponse struct {
	dig.In
}

func (BookControllerResponse) invalidMessage(ctx echo.Context, err error) error {
	return ctx.JSON(http.StatusBadRequest, GeneralResponse{
		Message: "Invalid Message",
	})
}

func (BookControllerResponse) invalidID(ctx echo.Context, err error) error {
	return ctx.JSON(http.StatusBadRequest, GeneralResponse{
		Message: "Invalid ID",
	})
}

func (BookControllerResponse) insertSuccess(ctx echo.Context, lastInsertID int64) error {
	return ctx.JSON(http.StatusCreated, GeneralResponse{
		Message: fmt.Sprintf("Success insert new record #%d", lastInsertID),
	})
}

func (BookControllerResponse) bookNotFound(ctx echo.Context, id int64) error {
	return ctx.JSON(http.StatusNotFound, GeneralResponse{
		Message: fmt.Sprintf("book #%d not found", id),
	})
}

func (BookControllerResponse) bookDeleted(ctx echo.Context, id int64) error {
	return ctx.JSON(http.StatusOK, GeneralResponse{
		Message: fmt.Sprintf("Delete #%d done", id),
	})
}

func (BookControllerResponse) bookUpdated(ctx echo.Context, id int64) error {
	return ctx.JSON(http.StatusOK, GeneralResponse{
		Message: fmt.Sprintf("Update #%d success", id),
	})
}
