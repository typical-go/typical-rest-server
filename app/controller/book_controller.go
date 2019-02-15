package controller

import (
	"github.com/labstack/echo"
)

// BookController handle input related to Book
type BookController interface {
	CRUD
}

type bookController struct {
}

// NewBookController return new instance of book controller
func NewBookController() BookController {
	return &bookController{}
}

func (c *bookController) Create(e echo.Context) error {
	return underContruction(e)
}

func (c *bookController) List(e echo.Context) error {
	return underContruction(e)
}

func (c *bookController) Get(e echo.Context) error {
	return underContruction(e)
}

func (c *bookController) Delete(e echo.Context) error {
	return underContruction(e)
}

func (c *bookController) Update(e echo.Context) error {
	return underContruction(e)
}
