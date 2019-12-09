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

// MusicCntrl is controller to Music entity
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
		return responsekit.InvalidRequest(ctx, err)
	}
	if lastInsertID, err = c.MusicService.Insert(ctx0, music); err != nil {
		return err
	}
	return responsekit.InsertSuccess(ctx, lastInsertID)
}

// List of music
func (c *MusicCntrl) List(ctx echo.Context) (err error) {
	var musics []*repository.Music
	ctx0 := ctx.Request().Context()
	if musics, err = c.MusicService.List(ctx0); err != nil {
		return err
	}
	return ctx.JSON(http.StatusOK, musics)
}

// Get music
func (c *MusicCntrl) Get(ctx echo.Context) (err error) {
	var id int64
	var music *repository.Music
	ctx0 := ctx.Request().Context()
	if id, err = strconv.ParseInt(ctx.Param("id"), 10, 64); err != nil {
		return responsekit.InvalidID(ctx, err)
	}
	if music, err = c.MusicService.Find(ctx0, id); err != nil {
		return err
	}
	if music == nil {
		return responsekit.NotFound(ctx, id)
	}
	return ctx.JSON(http.StatusOK, music)
}

// Delete music
func (c *MusicCntrl) Delete(ctx echo.Context) (err error) {
	var id int64
	ctx0 := ctx.Request().Context()
	if id, err = strconv.ParseInt(ctx.Param("id"), 10, 64); err != nil {
		return responsekit.InvalidID(ctx, err)
	}
	if err = c.MusicService.Delete(ctx0, id); err != nil {
		return err
	}
	return responsekit.DeleteSuccess(ctx, id)
}

// Update music
func (c *MusicCntrl) Update(ctx echo.Context) (err error) {
	var music repository.Music
	ctx0 := ctx.Request().Context()
	if err = ctx.Bind(&music); err != nil {
		return err
	}
	if music.ID <= 0 {
		return responsekit.InvalidID(ctx, err)
	}
	if err = validator.New().Struct(music); err != nil {
		return responsekit.InvalidRequest(ctx, err)
	}
	if err = c.MusicService.Update(ctx0, music); err != nil {
		return err
	}
	return responsekit.UpdateSuccess(ctx, music.ID)
}
