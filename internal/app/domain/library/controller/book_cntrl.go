package controller

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/typical-go/typical-rest-server/internal/app/data_access/librarydb"
	"github.com/typical-go/typical-rest-server/internal/app/domain/library/service"
	"github.com/typical-go/typical-rest-server/pkg/typrest"
	"go.uber.org/dig"
)

type (
	// BookCntrl is controller to book entity
	BookCntrl struct {
		dig.In
		service.BookSvc
	}
)

var _ typrest.Router = (*BookCntrl)(nil)

// SetRoute to define API Route
func (c *BookCntrl) SetRoute(e typrest.Server) {
	e.GET("/books", c.Find)
	e.GET("/books/:id", c.FindOne)
	e.POST("/books", c.Create)
	e.PUT("/books/:id", c.Update)
	e.PATCH("/books/:id", c.Patch)
	e.DELETE("/books/:id", c.Delete)
}

// Create book
func (c *BookCntrl) Create(ec echo.Context) (err error) {
	var book librarydb.Book
	if err = ec.Bind(&book); err != nil {
		return err
	}
	ctx := ec.Request().Context()
	newBook, err := c.BookSvc.Create(ctx, &book)
	if err != nil {
		return typrest.HTTPError(err)
	}
	ec.Response().Header().Set(echo.HeaderLocation, fmt.Sprintf("/books/%d", newBook.ID))
	return ec.JSON(http.StatusCreated, newBook)
}

// Find books
func (c *BookCntrl) Find(ec echo.Context) (err error) {
	var books []*librarydb.Book
	if books, err = c.BookSvc.Find(
		ec.Request().Context(),
	); err != nil {
		return typrest.HTTPError(err)
	}
	return ec.JSON(http.StatusOK, books)
}

// FindOne book
func (c *BookCntrl) FindOne(ec echo.Context) error {
	book, err := c.BookSvc.FindOne(
		ec.Request().Context(),
		ec.Param("id"),
	)
	if err != nil {
		return typrest.HTTPError(err)
	}
	return ec.JSON(http.StatusOK, book)
}

// Delete book
func (c *BookCntrl) Delete(ec echo.Context) (err error) {
	if err = c.BookSvc.Delete(
		ec.Request().Context(),
		ec.Param("id"),
	); err != nil {
		return typrest.HTTPError(err)
	}
	return ec.NoContent(http.StatusNoContent)
}

// Update book
func (c *BookCntrl) Update(ec echo.Context) (err error) {
	var book librarydb.Book
	if err = ec.Bind(&book); err != nil {
		return err
	}
	ctx := ec.Request().Context()
	paramID := ec.Param("id")
	updatedBook, err := c.BookSvc.Update(ctx, paramID, &book)
	if err != nil {
		return typrest.HTTPError(err)
	}
	return ec.JSON(http.StatusOK, updatedBook)
}

// Patch book
func (c *BookCntrl) Patch(ec echo.Context) (err error) {
	var book librarydb.Book
	if err = ec.Bind(&book); err != nil {
		return err
	}
	ctx := ec.Request().Context()
	paramID := ec.Param("id")
	patchedBook, err := c.BookSvc.Patch(ctx, paramID, &book)
	if err != nil {
		return typrest.HTTPError(err)
	}
	return ec.JSON(http.StatusOK, patchedBook)
}
