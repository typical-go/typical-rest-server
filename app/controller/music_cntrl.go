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
	e.GET("musics", c.Find)
	e.POST("musics", c.Create)
	e.GET("musics/:id", c.FindOne)
	e.PUT("musics", c.Update)
	e.DELETE("musics/:id", c.Delete)
}

// Create music
func (c *MusicCntrl) Create(ec echo.Context) (err error) {
	var (
		music        repository.Music
		lastInsertID int64
		ctx          = ec.Request().Context()
	)
	if err = ec.Bind(&music); err != nil {
		return err
	}
	if err = validator.New().Struct(music); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	if lastInsertID, err = c.MusicService.Create(ctx, music); err != nil {
		return echo.NewHTTPError(http.StatusUnprocessableEntity, err.Error())
	}
	return ec.JSON(http.StatusCreated, GeneralResponse{
		Message: fmt.Sprintf("Success insert new music #%d", lastInsertID),
	})
}

// Find musics
func (c *MusicCntrl) Find(ec echo.Context) (err error) {
	var (
		musics []*repository.Music
		ctx    = ec.Request().Context()
	)
	if musics, err = c.MusicService.Find(ctx); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return ec.JSON(http.StatusOK, musics)
}

// FindOne music
func (c *MusicCntrl) FindOne(ec echo.Context) (err error) {
	var (
		id    int64
		music *repository.Music
		ctx   = ec.Request().Context()
	)
	if id, err = strconv.ParseInt(ec.Param("id"), 10, 64); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid ID")
	}
	if music, err = c.MusicService.FindOne(ctx, id); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	if music == nil {
		return echo.NewHTTPError(http.StatusNotFound, fmt.Sprintf("Music#%d not found", id))
	}
	return ec.JSON(http.StatusOK, music)
}

// Delete music
func (c *MusicCntrl) Delete(ec echo.Context) (err error) {
	var (
		id   int64
		ctx0 = ec.Request().Context()
	)
	if id, err = strconv.ParseInt(ec.Param("id"), 10, 64); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid ID")
	}
	if err = c.MusicService.Delete(ctx0, id); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return ec.JSON(http.StatusOK, GeneralResponse{
		Message: fmt.Sprintf("Success delete music #%d", id),
	})
}

// Update music
func (c *MusicCntrl) Update(ec echo.Context) (err error) {
	var (
		music repository.Music
		ctx   = ec.Request().Context()
	)
	if err = ec.Bind(&music); err != nil {
		return err
	}
	if music.ID <= 0 {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid ID")
	}
	if err = validator.New().Struct(music); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	if err = c.MusicService.Update(ctx, music); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return ec.JSON(http.StatusOK, GeneralResponse{
		Message: fmt.Sprintf("Success update music #%d", music.ID),
	})
}
