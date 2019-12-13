package controller

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/labstack/echo"

	"github.com/typical-go/typical-rest-server/app/repository"
	"github.com/typical-go/typical-rest-server/app/service"
	"go.uber.org/dig"
	"gopkg.in/go-playground/validator.v9"
)

// BookCntrl is controller to book entity
type BookCntrl struct {
	dig.In
	service.BookService
}

// Route to define API Route
func (c *BookCntrl) Route(e *echo.Echo) {
	e.GET("books", c.List)
	e.POST("books", c.Create)
	e.GET("books/:id", c.Get)
	e.PUT("books", c.Update)
	e.DELETE("books/:id", c.Delete)
}

// Create book
func (c *BookCntrl) Create(ctx echo.Context) (err error) {
	var book repository.Book
	var lastInsertID int64
	ctx0 := ctx.Request().Context()
	if err = ctx.Bind(&book); err != nil {
		return err
	}
	if err = validator.New().Struct(book); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	if lastInsertID, err = c.BookService.Insert(ctx0, book); err != nil {
		return echo.NewHTTPError(http.StatusUnprocessableEntity, err.Error())
	}
	return ctx.JSON(http.StatusCreated, GeneralResponse{
		Message: fmt.Sprintf("Success insert new book #%d", lastInsertID),
	})
}

// List of book
func (c *BookCntrl) List(ctx echo.Context) (err error) {
	var books []*repository.Book
	ctx0 := ctx.Request().Context()
	if books, err = c.BookService.List(ctx0); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return ctx.JSON(http.StatusOK, books)
}

// Get book
func (c *BookCntrl) Get(ctx echo.Context) (err error) {
	var id int64
	var book *repository.Book
	ctx0 := ctx.Request().Context()
	if id, err = strconv.ParseInt(ctx.Param("id"), 10, 64); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid ID")
	}
	if book, err = c.BookService.Find(ctx0, id); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	if book == nil {
		return echo.NewHTTPError(http.StatusNotFound, fmt.Sprintf("Book#%d not found", id))
	}
	return ctx.JSON(http.StatusOK, book)
}

// Delete book
func (c *BookCntrl) Delete(ctx echo.Context) (err error) {
	var id int64
	ctx0 := ctx.Request().Context()
	if id, err = strconv.ParseInt(ctx.Param("id"), 10, 64); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid ID")
	}
	if err = c.BookService.Delete(ctx0, id); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return ctx.JSON(http.StatusOK, GeneralResponse{
		Message: fmt.Sprintf("Success delete book #%d", id),
	})
}

// Update book
func (c *BookCntrl) Update(ctx echo.Context) (err error) {
	var book repository.Book
	ctx0 := ctx.Request().Context()
	if err = ctx.Bind(&book); err != nil {
		return err
	}
	if book.ID <= 0 {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid ID")
	}
	if err = validator.New().Struct(book); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	if err = c.BookService.Update(ctx0, book); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return ctx.JSON(http.StatusOK, GeneralResponse{
		Message: fmt.Sprintf("Success update book #%d", book.ID),
	})
}
