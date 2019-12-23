package controller

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/labstack/echo"

	"github.com/typical-go/typical-rest-server/app/repository"
	"github.com/typical-go/typical-rest-server/app/service"
	"go.uber.org/dig"
	"gopkg.in/go-playground/validator.v9"
)

// LocaleCntrl is controller to locale entity
type LocaleCntrl struct {
	dig.In
	service.LocaleService
}

// Route to define API Route
func (c *LocaleCntrl) Route(e *echo.Echo) {
	e.GET("locales", c.List)
	e.POST("locales", c.Create)
	e.GET("locales/:id", c.Get)
	e.PUT("locales", c.Update)
	e.DELETE("locales/:id", c.Delete)
}

// Create locale
func (c *LocaleCntrl) Create(ctx echo.Context) (err error) {
	var locale repository.Locale
	var lastInsertID int64
	ctx0 := ctx.Request().Context()
	if err = ctx.Bind(&locale); err != nil {
		return err
	}
	if err = validator.New().Struct(locale); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	if lastInsertID, err = c.LocaleService.Insert(ctx0, locale); err != nil {
		return echo.NewHTTPError(http.StatusUnprocessableEntity, err.Error())
	}
	return ctx.JSON(http.StatusCreated, GeneralResponse{
		Message: fmt.Sprintf("Success insert new locale #%d", lastInsertID),
	})
}

// List of locale
func (c *LocaleCntrl) List(ctx echo.Context) (err error) {
	var locales []*repository.Locale
	ctx0 := ctx.Request().Context()
	if locales, err = c.LocaleService.List(ctx0); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return ctx.JSON(http.StatusOK, locales)
}

// Get locale
func (c *LocaleCntrl) Get(ctx echo.Context) (err error) {
	var id int64
	var locale *repository.Locale
	ctx0 := ctx.Request().Context()
	if id, err = strconv.ParseInt(ctx.Param("id"), 10, 64); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid ID")
	}
	if locale, err = c.LocaleService.Find(ctx0, id); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	if locale == nil {
		return echo.NewHTTPError(http.StatusNotFound, fmt.Sprintf("Locale#%d not found", id))
	}
	return ctx.JSON(http.StatusOK, locale)
}

// Delete locale
func (c *LocaleCntrl) Delete(ctx echo.Context) (err error) {
	var id int64
	ctx0 := ctx.Request().Context()
	if id, err = strconv.ParseInt(ctx.Param("id"), 10, 64); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid ID")
	}
	if err = c.LocaleService.Delete(ctx0, id); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return ctx.JSON(http.StatusOK, GeneralResponse{
		Message: fmt.Sprintf("Success delete locale #%d", id),
	})
}

// Update locale
func (c *LocaleCntrl) Update(ctx echo.Context) (err error) {
	var locale repository.Locale
	ctx0 := ctx.Request().Context()
	if err = ctx.Bind(&locale); err != nil {
		return err
	}
	if locale.ID <= 0 {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid ID")
	}
	if err = validator.New().Struct(locale); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	if err = c.LocaleService.Update(ctx0, locale); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return ctx.JSON(http.StatusOK, GeneralResponse{
		Message: fmt.Sprintf("Success update locale #%d", locale.ID),
	})
}
