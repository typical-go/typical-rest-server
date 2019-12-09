package controller

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo"

	"github.com/typical-go/typical-rest-server/app/repository"
	"github.com/typical-go/typical-rest-server/app/service"
	"github.com/typical-go/typical-rest-server/pkg/utility/responsekit"
	"go.uber.org/dig"
	"gopkg.in/go-playground/validator.v9"
)

// BookCntrl is controller to book entity
type MusicCntrl struct {
	dig.In
	service.BookService
}

// Route to define API Route
func (c *MusicCntrl) Route(e *echo.Echo) {
	e.GET("musics", c.List)
	e.POST("musics", c.Create)
	e.GET("musics/:id", c.Get)
	e.PUT("musics", c.Update)
	e.DELETE("musics/:id", c.Delete)
}

// Create book
func (c *MusicCntrl) Create(ctx echo.Context) (err error) {
	var book repository.Book
	var lastInsertID int64
	ctx0 := ctx.Request().Context()
	if err = ctx.Bind(&book); err != nil {
		return err
	}
	if err = validator.New().Struct(book); err != nil {
		return responsekit.InvalidRequest(ctx, err)
	}
	if lastInsertID, err = c.BookService.Insert(ctx0, book); err != nil {
		return err
	}
	return responsekit.InsertSuccess(ctx, lastInsertID)
}

// List of book
func (c *MusicCntrl) List(ctx echo.Context) (err error) {
	var books []*repository.Book
	ctx0 := ctx.Request().Context()
	if books, err = c.BookService.List(ctx0); err != nil {
		return err
	}
	return ctx.JSON(http.StatusOK, books)
}

// Get book
func (c *MusicCntrl) Get(ctx echo.Context) (err error) {
	var id int64
	var book *repository.Book
	ctx0 := ctx.Request().Context()
	if id, err = strconv.ParseInt(ctx.Param("id"), 10, 64); err != nil {
		return responsekit.InvalidID(ctx, err)
	}
	if book, err = c.BookService.Find(ctx0, id); err != nil {
		return err
	}
	if book == nil {
		return responsekit.NotFound(ctx, id)
	}
	return ctx.JSON(http.StatusOK, book)
}

// Delete book
func (c *MusicCntrl) Delete(ctx echo.Context) (err error) {
	var id int64
	ctx0 := ctx.Request().Context()
	if id, err = strconv.ParseInt(ctx.Param("id"), 10, 64); err != nil {
		return responsekit.InvalidID(ctx, err)
	}
	if err = c.BookService.Delete(ctx0, id); err != nil {
		return err
	}
	return responsekit.DeleteSuccess(ctx, id)
}

// Update book
func (c *MusicCntrl) Update(ctx echo.Context) (err error) {
	var book repository.Book
	ctx0 := ctx.Request().Context()
	if err = ctx.Bind(&book); err != nil {
		return err
	}
	if book.ID <= 0 {
		return responsekit.InvalidID(ctx, err)
	}
	if err = validator.New().Struct(book); err != nil {
		return responsekit.InvalidRequest(ctx, err)
	}
	if err = c.BookService.Update(ctx0, book); err != nil {
		return err
	}
	return responsekit.UpdateSuccess(ctx, book.ID)
}
