package controller

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/typical-go/typical-rest-server/internal/app/data_access/postgresdb"
	"github.com/typical-go/typical-rest-server/internal/app/domain/mylibrary/service"
	"github.com/typical-go/typical-rest-server/pkg/cachekit"
	"github.com/typical-go/typical-rest-server/pkg/typrest"
	"go.uber.org/dig"
)

type (
	// BookCntrl is controller to book entity
	BookCntrl struct {
		dig.In
		Svc   service.BookSvc
		Cache *cachekit.Store
	}
)

var _ typrest.Router = (*BookCntrl)(nil)

// SetRoute to define API Route
func (c *BookCntrl) SetRoute(e typrest.Server) {
	e.GET("/books", c.Find, c.Cache.Middleware)
	e.GET("/books/:id", c.FindOne, c.Cache.Middleware)
	e.HEAD("/books/:id", c.FindOne, c.Cache.Middleware)
	e.POST("/books", c.Create)
	e.PUT("/books/:id", c.Update)
	e.PATCH("/books/:id", c.Patch)
	e.DELETE("/books/:id", c.Delete)
}

// Create book
func (c *BookCntrl) Create(ec echo.Context) (err error) {
	var book postgresdb.Book
	if err = ec.Bind(&book); err != nil {
		return err
	}
	ctx := ec.Request().Context()
	newBook, err := c.Svc.Create(ctx, &book)
	if err != nil {
		return typrest.HTTPError(err)
	}
	ec.Response().Header().Set(echo.HeaderLocation, fmt.Sprintf("/books/%d", newBook.ID))
	return ec.JSON(http.StatusCreated, newBook)
}

// Find books
func (c *BookCntrl) Find(ec echo.Context) (err error) {
	var req service.FindReq
	if err = ec.Bind(&req); err != nil {
		return err
	}
	ctx := ec.Request().Context()
	books, err := c.Svc.Find(ctx, &req)
	if err != nil {
		return typrest.HTTPError(err)
	}
	return ec.JSON(http.StatusOK, books)
}

// FindOne book
func (c *BookCntrl) FindOne(ec echo.Context) error {
	book, err := c.Svc.FindOne(
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
	ctx := ec.Request().Context()
	id := ec.Param("id")
	if err = c.Svc.Delete(ctx, id); err != nil {
		return typrest.HTTPError(err)
	}
	return ec.NoContent(http.StatusNoContent)
}

// Update book
func (c *BookCntrl) Update(ec echo.Context) (err error) {
	var book postgresdb.Book
	if err = ec.Bind(&book); err != nil {
		return err
	}
	ctx := ec.Request().Context()
	paramID := ec.Param("id")
	updatedBook, err := c.Svc.Update(ctx, paramID, &book)
	if err != nil {
		return typrest.HTTPError(err)
	}
	return ec.JSON(http.StatusOK, updatedBook)
}

// Patch book
func (c *BookCntrl) Patch(ec echo.Context) (err error) {
	var book postgresdb.Book
	if err = ec.Bind(&book); err != nil {
		return err
	}
	ctx := ec.Request().Context()
	paramID := ec.Param("id")
	patchedBook, err := c.Svc.Patch(ctx, paramID, &book)
	if err != nil {
		return typrest.HTTPError(err)
	}
	return ec.JSON(http.StatusOK, patchedBook)
}
