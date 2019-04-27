package cntrl

import (
	"fmt"
	"net/http"

	"github.com/typical-go/typical-rest-server/app/helper/strkit"
	"github.com/typical-go/typical-rest-server/app/repo"
	"github.com/labstack/echo"
)

// BookController handle input related to Book
type BookController interface {
	CRUDController
}

type bookController struct {
	bookRepository repo.BookRepository
}

// NewBookController return new instance of book controller
func NewBookController(bookRepository repo.BookRepository) BookController {
	return &bookController{
		bookRepository: bookRepository,
	}
}

func (c *bookController) Create(ctx echo.Context) (err error) {
	if err := c.Check(); err != nil {
		return err
	}

	var book repo.Book

	err = ctx.Bind(&book)
	if err != nil {
		return err
	}

	err = book.Validate()
	if err != nil {
		return invalidMessage(ctx, err)
	}

	result, err := c.bookRepository.Insert(book)
	if err != nil {
		return err
	}

	return insertSuccess(ctx, result)

}

func (c *bookController) List(ctx echo.Context) error {
	if err := c.Check(); err != nil {
		return err
	}

	books, err := c.bookRepository.List()
	if err != nil {
		return err
	}
	return ctx.JSON(http.StatusOK, books)
}

func (c *bookController) Get(ctx echo.Context) error {
	if err := c.Check(); err != nil {
		return err
	}

	id, err := strkit.ToInt64(ctx.Param("id"))
	if err != nil {
		return invalidID(ctx, err)
	}

	book, err := c.bookRepository.Get(id)
	if err != nil {
		return err
	}

	if book == nil {
		return ctx.JSON(http.StatusNotFound, map[string]string{"message": fmt.Sprintf("book #%d not found", id)})
	}

	return ctx.JSON(http.StatusOK, book)
}

func (c *bookController) Delete(ctx echo.Context) error {
	if err := c.Check(); err != nil {
		return err
	}

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

func (c *bookController) Update(ctx echo.Context) (err error) {
	if err := c.Check(); err != nil {
		return err
	}

	var book repo.Book

	err = ctx.Bind(&book)
	if err != nil {
		return err
	}

	if book.ID <= 0 {
		return invalidID(ctx, err)
	}

	err = book.Validate()
	if err != nil {
		return invalidMessage(ctx, err)
	}

	err = c.bookRepository.Update(book)
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, map[string]string{"message": "Update success"})
}

func (c *bookController) Check() error {
	if c.bookRepository == nil {
		return fmt.Errorf("BookRepository is missing")
	}
	return nil
}
