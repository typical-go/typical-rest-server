package controller

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/imantung/typical-go-server/app/repository"
	"github.com/labstack/echo"
)

// BookController handle input related to Book
type BookController interface {
	CRUD
}

type bookController struct {
	bookRepository repository.BookRepository
}

// NewBookController return new instance of book controller
func NewBookController(bookRepository repository.BookRepository) BookController {
	return &bookController{
		bookRepository: bookRepository,
	}
}

func (c *bookController) Create(ctx echo.Context) (err error) {
	var book repository.Book

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
	books, err := c.bookRepository.List()
	if err != nil {
		return err
	}
	return ctx.JSON(http.StatusOK, books)
}

func (c *bookController) Get(ctx echo.Context) error {
	id, _ := strconv.ParseInt(ctx.Param("id"), 10, 64)
	book, err := c.bookRepository.Get(id)
	if err != nil {
		return err
	}
	return ctx.JSON(http.StatusOK, book)
}

func (c *bookController) Delete(ctx echo.Context) error {
	return underContruction(ctx)
}

func (c *bookController) Update(ctx echo.Context) error {
	return underContruction(ctx)
}

func (c *bookController) BeforeActionFunc(next echo.HandlerFunc) echo.HandlerFunc {
	return func(e echo.Context) error {
		err := c.beforeAction()
		if err != nil {
			return err
		}
		return next(e)
	}
}

func (c *bookController) beforeAction() error {
	if c.bookRepository == nil {
		return fmt.Errorf("BookRepository is missing")
	}
	return nil
}

func (c *bookController) RegisterTo(entity string, e *echo.Echo) {
	e.GET(fmt.Sprintf("/%s", entity), c.List, c.BeforeActionFunc)
	e.POST(fmt.Sprintf("/%s", entity), c.Create, c.BeforeActionFunc)
	e.GET(fmt.Sprintf("/%s/:id", entity), c.Get, c.BeforeActionFunc)
	e.PUT(fmt.Sprintf("/%s/:id", entity), c.Update, c.BeforeActionFunc)
	e.DELETE(fmt.Sprintf("/%s/:id", entity), c.Delete, c.BeforeActionFunc)
}
