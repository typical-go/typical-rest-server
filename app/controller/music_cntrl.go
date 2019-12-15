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

// MusicCntrl is controller to music entity
type MusicCntrl struct {
	dig.In
	service.MusicService
}

// Route to define API Route
func (c *MusicCntrl) Route(e *echo.Echo) {
	e.GET("musics", c.List)
	e.POST("musics", c.Create)
	e.GET("musics/:id", c.Get)
	e.PUT("musics", c.Update)
	e.DELETE("musics/:id", c.Delete)
}

// Create music
func (c *MusicCntrl) Create(ctx echo.Context) (err error) {
	var music repository.Music
	var lastInsertID int64
	ctx0 := ctx.Request().Context()
	if err = ctx.Bind(&music); err != nil {
		return err
	}
	if err = validator.New().Struct(music); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	if lastInsertID, err = c.MusicService.Insert(ctx0, music); err != nil {
		return echo.NewHTTPError(http.StatusUnprocessableEntity, err.Error())
	}
	return ctx.JSON(http.StatusCreated, GeneralResponse{
		Message: fmt.Sprintf("Success insert new music #%d", lastInsertID),
	})
}

// List of music
func (c *MusicCntrl) List(ctx echo.Context) (err error) {
	var musics []*repository.Music
	ctx0 := ctx.Request().Context()
	if musics, err = c.MusicService.List(ctx0); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return ctx.JSON(http.StatusOK, musics)
}

// Get music
func (c *MusicCntrl) Get(ctx echo.Context) (err error) {
	var id int64
	var music *repository.Music
	ctx0 := ctx.Request().Context()
	if id, err = strconv.ParseInt(ctx.Param("id"), 10, 64); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid ID")
	}
	if music, err = c.MusicService.Find(ctx0, id); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	if music == nil {
		return echo.NewHTTPError(http.StatusNotFound, fmt.Sprintf("Music#%d not found", id))
	}
	return ctx.JSON(http.StatusOK, music)
}

// Delete music
func (c *MusicCntrl) Delete(ctx echo.Context) (err error) {
	var id int64
	ctx0 := ctx.Request().Context()
	if id, err = strconv.ParseInt(ctx.Param("id"), 10, 64); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid ID")
	}
	if err = c.MusicService.Delete(ctx0, id); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return ctx.JSON(http.StatusOK, GeneralResponse{
		Message: fmt.Sprintf("Success delete music #%d", id),
	})
}

// Update music
func (c *MusicCntrl) Update(ctx echo.Context) (err error) {
	var music repository.Music
	ctx0 := ctx.Request().Context()
	if err = ctx.Bind(&music); err != nil {
		return err
	}
	if music.ID <= 0 {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid ID")
	}
	if err = validator.New().Struct(music); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	if err = c.MusicService.Update(ctx0, music); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return ctx.JSON(http.StatusOK, GeneralResponse{
		Message: fmt.Sprintf("Success update music #%d", music.ID),
	})
}
