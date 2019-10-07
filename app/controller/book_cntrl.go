package controller

import (
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
	BookCntrlResponse
	service.BookService
}

// Route to define API Route
func (c *BookCntrl) Route(e *echo.Echo) {
	e.GET("book", c.List)
	e.POST("book", c.Create)
	e.GET("book/:id", c.Get)
	e.PUT("book", c.Update)
	e.DELETE("book/:id", c.Delete)
}

// Create book
func (c *BookCntrl) Create(ctx echo.Context) (err error) {
	var book repository.Book
	err = ctx.Bind(&book)
	if err != nil {
		return err
	}
	err = validator.New().Struct(book)
	if err != nil {
		return c.invalidMessage(ctx, err)
	}
	ctx0 := ctx.Request().Context()
	result, err := c.BookService.Insert(ctx0, book)
	if err != nil {
		return err
	}
	return c.insertSuccess(ctx, result)
}

// List of book
func (c *BookCntrl) List(ctx echo.Context) error {
	ctx0 := ctx.Request().Context()
	books, err := c.BookService.List(ctx0)
	if err != nil {
		return err
	}
	return ctx.JSON(http.StatusOK, books)
}

// Get book
func (c *BookCntrl) Get(ctx echo.Context) error {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		return c.invalidID(ctx, err)
	}
	ctx0 := ctx.Request().Context()
	book, err := c.BookService.Find(ctx0, id)
	if err != nil {
		return err
	}
	if book == nil {
		return c.bookNotFound(ctx, id)
	}
	return ctx.JSON(http.StatusOK, book)
}

// Delete book
func (c *BookCntrl) Delete(ctx echo.Context) error {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		return c.invalidID(ctx, err)
	}
	ctx0 := ctx.Request().Context()
	err = c.BookService.Delete(ctx0, id)
	if err != nil {
		return err
	}
	return c.bookDeleted(ctx, id)
}

// Update book
func (c *BookCntrl) Update(ctx echo.Context) (err error) {
	var book repository.Book
	err = ctx.Bind(&book)
	if err != nil {
		return err
	}
	if book.ID <= 0 {
		return c.invalidID(ctx, err)
	}
	err = validator.New().Struct(book)
	if err != nil {
		return c.invalidMessage(ctx, err)
	}
	ctx0 := ctx.Request().Context()
	err = c.BookService.Update(ctx0, book)
	if err != nil {
		return err
	}
	return c.bookUpdated(ctx, book.ID)
}
