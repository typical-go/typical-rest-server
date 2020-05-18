package controller

import (
	"net/http"

	"github.com/labstack/echo"

	"github.com/typical-go/typical-rest-server/internal/server/repository"
	"github.com/typical-go/typical-rest-server/internal/server/service"
	"go.uber.org/dig"
)

type (

	// BookCntrl is controller to book entity
	BookCntrl struct {
		dig.In
		service.BookService
	}
)

// SetRoute to define API Route
func (c *BookCntrl) SetRoute(e *echo.Echo) {
	e.GET("books", c.Find)
	e.POST("books", c.Create)
	e.GET("books/:id", c.FindOne)
	e.PUT("books/:id", c.Update)
	e.DELETE("books/:id", c.Delete)
}

// Create book
func (c *BookCntrl) Create(ec echo.Context) (err error) {
	var (
		result *repository.Book
		book   repository.Book
	)

	if err = ec.Bind(&book); err != nil {
		return err
	}

	if result, err = c.BookService.Create(
		ec.Request().Context(),
		&book,
	); err != nil {
		return httpError(err)
	}

	return ec.JSON(http.StatusCreated, result)
}

// Find books
func (c *BookCntrl) Find(ec echo.Context) (err error) {
	var books []*repository.Book
	if books, err = c.BookService.Find(
		ec.Request().Context(),
	); err != nil {
		return httpError(err)
	}
	return ec.JSON(http.StatusOK, books)
}

// FindOne book
func (c *BookCntrl) FindOne(ec echo.Context) error {
	book, err := c.BookService.FindOne(
		ec.Request().Context(),
		ec.Param("id"),
	)

	if err != nil {
		return httpError(err)
	}

	return ec.JSON(http.StatusOK, book)
}

// Delete book
func (c *BookCntrl) Delete(ec echo.Context) (err error) {
	if err = c.BookService.Delete(
		ec.Request().Context(),
		ec.Param("id"),
	); err != nil {
		return httpError(err)
	}
	return ec.NoContent(http.StatusNoContent)
}

// Update book
func (c *BookCntrl) Update(ec echo.Context) error {
	var (
		form   repository.Book
		result *repository.Book
		err    error
	)

	if err = ec.Bind(&form); err != nil {
		return err
	}

	if result, err = c.BookService.Update(
		ec.Request().Context(),
		ec.Param("id"),
		&form,
	); err != nil {
		return httpError(err)
	}

	return ec.JSON(http.StatusOK, result)
}
