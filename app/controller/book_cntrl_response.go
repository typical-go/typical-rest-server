package controller

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo"
	"go.uber.org/dig"
)

// BookCntrlResponse is response of BookCntrl
type BookCntrlResponse struct {
	dig.In
}

func (BookCntrlResponse) invalidMessage(ctx echo.Context, err error) error {
	return ctx.JSON(http.StatusBadRequest, GeneralResponse{
		Message: "Invalid Message",
	})
}

func (BookCntrlResponse) invalidID(ctx echo.Context, err error) error {
	return ctx.JSON(http.StatusBadRequest, GeneralResponse{
		Message: "Invalid ID",
	})
}

func (BookCntrlResponse) insertSuccess(ctx echo.Context, lastInsertID int64) error {
	return ctx.JSON(http.StatusCreated, GeneralResponse{
		Message: fmt.Sprintf("Success insert new record #%d", lastInsertID),
	})
}

func (BookCntrlResponse) bookNotFound(ctx echo.Context, id int64) error {
	return ctx.JSON(http.StatusNotFound, GeneralResponse{
		Message: fmt.Sprintf("book #%d not found", id),
	})
}

func (BookCntrlResponse) bookDeleted(ctx echo.Context, id int64) error {
	return ctx.JSON(http.StatusOK, GeneralResponse{
		Message: fmt.Sprintf("Delete #%d done", id),
	})
}

func (BookCntrlResponse) bookUpdated(ctx echo.Context, id int64) error {
	return ctx.JSON(http.StatusOK, GeneralResponse{
		Message: fmt.Sprintf("Update #%d success", id),
	})
}
