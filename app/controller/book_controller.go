package controller

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo"
	"github.com/typical-go/typical-rest-server/app/repository"
	"github.com/typical-go/typical-rest-server/pkg/utility/strkit"
	"gopkg.in/go-playground/validator.v9"
)

// BookController is controller to book entity
type BookController struct {
	bookRepository repository.BookRepository
}

// NewBookController return new instance of book controller
func NewBookController(bookRepository repository.BookRepository) *BookController {
	return &BookController{
		bookRepository: bookRepository,
	}
}

// Route to define API Route
func (c *BookController) Route(e *echo.Echo) {
	e.GET("book", c.List)
	e.POST("book", c.Create)
	e.GET("book/:id", c.Get)
	e.PUT("book", c.Update)
	e.DELETE("book/:id", c.Delete)
}

// Create book
func (c *BookController) Create(ctx echo.Context) (err error) {
	var book repository.Book
	err = ctx.Bind(&book)
	if err != nil {
		return err
	}
	err = validator.New().Struct(book)
	if err != nil {
		return invalidMessage(ctx, err)
	}
	result, err := c.bookRepository.Insert(book)
	if err != nil {
		return err
	}
	return insertSuccess(ctx, result)
}

// List of book
func (c *BookController) List(ctx echo.Context) error {
	books, err := c.bookRepository.List()
	if err != nil {
		return err
	}
	return ctx.JSON(http.StatusOK, books)
}

// Get book
func (c *BookController) Get(ctx echo.Context) error {
	id, err := strkit.ToInt64(ctx.Param("id"))
	if err != nil {
		return invalidID(ctx, err)
	}
	book, err := c.bookRepository.Find(id)
	if err != nil {
		return err
	}
	if book == nil {
		return ctx.JSON(http.StatusNotFound, map[string]string{"message": fmt.Sprintf("book #%d not found", id)})
	}
	return ctx.JSON(http.StatusOK, book)
}

// Delete book
func (c *BookController) Delete(ctx echo.Context) error {
	id, err := strkit.ToInt64(ctx.Param("id"))
	if err != nil {
		return invalidID(ctx, err)
	}
	err = c.bookRepository.Delete(id)
	if err != nil {
		return err
	}
	return ctx.JSON(http.StatusOK, map[string]string{"message": fmt.Sprintf("Delete #%d done", id)})
}

// Update book
func (c *BookController) Update(ctx echo.Context) (err error) {
	var book repository.Book
	err = ctx.Bind(&book)
	if err != nil {
		return err
	}
	if book.ID <= 0 {
		return invalidID(ctx, err)
	}
	err = validator.New().Struct(book)
	if err != nil {
		return invalidMessage(ctx, err)
	}
	err = c.bookRepository.Update(book)
	if err != nil {
		return err
	}
	return ctx.JSON(http.StatusOK, map[string]string{"message": "Update success"})
}
