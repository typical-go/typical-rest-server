package controller

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/labstack/echo"

	"github.com/typical-go/typical-rest-server/internal/server/repository"
	"github.com/typical-go/typical-rest-server/internal/server/service"
	"github.com/typical-go/typical-rest-server/pkg/echokit"
	"github.com/typical-go/typical-rest-server/pkg/errvalid"
	"go.uber.org/dig"
)

type (
	// BookCntrl is controller to book entity
	BookCntrl struct {
		dig.In
		service.BookSvc
	}
)

var _ echokit.Router = (*BookCntrl)(nil)

// Route to define API Route
func (c *BookCntrl) Route(e echokit.Server) (err error) {
	e.GET("books", c.Retrieve)
	e.GET("books/:id", c.RetrieveOne)
	e.POST("books", c.Create)
	e.PUT("books/:id", c.Update)
	e.PATCH("books/:id", c.Patch)
	e.DELETE("books/:id", c.Delete)
	return
}

// Create book
func (c *BookCntrl) Create(ec echo.Context) (err error) {
	var book repository.Book
	if err = ec.Bind(&book); err != nil {
		return err
	}

	ctx := ec.Request().Context()
	newBook, err := c.BookSvc.Create(ctx, &book)
	if err != nil {
		return c.HandleError(err)
	}

	ec.Response().Header().Set(echo.HeaderLocation, fmt.Sprintf("/books/%d", newBook.ID))
	return ec.JSON(http.StatusCreated, newBook)
}

// Retrieve books
func (c *BookCntrl) Retrieve(ec echo.Context) (err error) {
	var books []*repository.Book
	if books, err = c.BookSvc.Retrieve(
		ec.Request().Context(),
	); err != nil {
		return c.HandleError(err)
	}
	return ec.JSON(http.StatusOK, books)
}

// RetrieveOne book
func (c *BookCntrl) RetrieveOne(ec echo.Context) error {
	book, err := c.BookSvc.RetrieveOne(
		ec.Request().Context(),
		ec.Param("id"),
	)

	if err != nil {
		return c.HandleError(err)
	}

	return ec.JSON(http.StatusOK, book)
}

// Delete book
func (c *BookCntrl) Delete(ec echo.Context) (err error) {
	if err = c.BookSvc.Delete(
		ec.Request().Context(),
		ec.Param("id"),
	); err != nil {
		return c.HandleError(err)
	}
	return ec.NoContent(http.StatusNoContent)
}

// Update book
func (c *BookCntrl) Update(ec echo.Context) (err error) {
	var book repository.Book
	if err = ec.Bind(&book); err != nil {
		return err
	}

	ctx := ec.Request().Context()
	paramID := ec.Param("id")
	updatedBook, err := c.BookSvc.Update(ctx, paramID, &book)
	if err != nil {
		return c.HandleError(err)
	}

	return ec.JSON(http.StatusOK, updatedBook)
}

// Patch book
func (c *BookCntrl) Patch(ec echo.Context) (err error) {
	var book repository.Book
	if err = ec.Bind(&book); err != nil {
		return err
	}

	ctx := ec.Request().Context()
	paramID := ec.Param("id")
	patchedBook, err := c.BookSvc.Patch(ctx, paramID, &book)
	if err != nil {
		return c.HandleError(err)
	}

	return ec.JSON(http.StatusOK, patchedBook)
}

// HandleError to handle error on this controller
func (c *BookCntrl) HandleError(err error) *echo.HTTPError {
	if err == sql.ErrNoRows {
		return echo.NewHTTPError(http.StatusNotFound)
	}

	if errvalid.Check(err) {
		return echo.NewHTTPError(
			http.StatusUnprocessableEntity,
			errvalid.Message(err),
		)
	}

	return echo.NewHTTPError(
		http.StatusInternalServerError,
		err.Error(),
	)
}
