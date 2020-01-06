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
	e.GET("books", c.Find)
	e.POST("books", c.Create)
	e.GET("books/:id", c.FindOne)
	e.PUT("books", c.Update)
	e.DELETE("books/:id", c.Delete)
}

// Create book
func (c *BookCntrl) Create(ec echo.Context) (err error) {
	var (
		book         repository.Book
		lastInsertID int64
		ctx          = ec.Request().Context()
	)
	if err = ec.Bind(&book); err != nil {
		return err
	}
	if err = validator.New().Struct(book); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	if lastInsertID, err = c.BookService.Create(ctx, book); err != nil {
		return echo.NewHTTPError(http.StatusUnprocessableEntity, err.Error())
	}
	return ec.JSON(http.StatusCreated, GeneralResponse{
		Message: fmt.Sprintf("Success insert new book #%d", lastInsertID),
	})
}

// Find books
func (c *BookCntrl) Find(ec echo.Context) (err error) {
	var (
		books []*repository.Book
		ctx   = ec.Request().Context()
	)
	if books, err = c.BookService.Find(ctx); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return ec.JSON(http.StatusOK, books)
}

// FindOne book
func (c *BookCntrl) FindOne(ec echo.Context) (err error) {
	var (
		id   int64
		book *repository.Book
		ctx  = ec.Request().Context()
	)
	if id, err = strconv.ParseInt(ec.Param("id"), 10, 64); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid ID")
	}
	if book, err = c.BookService.FindOne(ctx, id); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	if book == nil {
		return echo.NewHTTPError(http.StatusNotFound, fmt.Sprintf("Book#%d not found", id))
	}
	return ec.JSON(http.StatusOK, book)
}

// Delete book
func (c *BookCntrl) Delete(ec echo.Context) (err error) {
	var (
		id  int64
		ctx = ec.Request().Context()
	)
	if id, err = strconv.ParseInt(ec.Param("id"), 10, 64); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid ID")
	}
	if err = c.BookService.Delete(ctx, id); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return ec.JSON(http.StatusOK, GeneralResponse{
		Message: fmt.Sprintf("Success delete book #%d", id),
	})
}

// Update book
func (c *BookCntrl) Update(ec echo.Context) (err error) {
	var (
		book repository.Book
		ctx  = ec.Request().Context()
	)
	if err = ec.Bind(&book); err != nil {
		return err
	}
	if book.ID <= 0 {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid ID")
	}
	if err = validator.New().Struct(book); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	if err = c.BookService.Update(ctx, book); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return ec.JSON(http.StatusOK, GeneralResponse{
		Message: fmt.Sprintf("Success update book #%d", book.ID),
	})
}
