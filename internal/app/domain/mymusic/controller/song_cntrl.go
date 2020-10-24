package controller

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/typical-go/typical-rest-server/internal/app/data_access/mysqldb"
	"github.com/typical-go/typical-rest-server/internal/app/domain/mymusic/service"
	"github.com/typical-go/typical-rest-server/pkg/typrest"
	"go.uber.org/dig"
)

type (
	// SongCntrl is controller to book entity
	SongCntrl struct {
		dig.In
		service.SongSvc
	}
)

var _ typrest.Router = (*SongCntrl)(nil)

// SetRoute to define API Route
func (c *SongCntrl) SetRoute(e typrest.Server) {
	e.GET("/songs", c.Find)
	e.GET("/songs/:id", c.FindOne)
	e.HEAD("/songs/:id", c.FindOne)
	e.POST("/songs", c.Create)
	e.PUT("/songs/:id", c.Update)
	e.PATCH("/songs/:id", c.Patch)
	e.DELETE("/songs/:id", c.Delete)
}

// Create book
func (c *SongCntrl) Create(ec echo.Context) (err error) {
	var book mysqldb.Song
	if err = ec.Bind(&book); err != nil {
		return err
	}
	ctx := ec.Request().Context()
	newSong, err := c.SongSvc.Create(ctx, &book)
	if err != nil {
		return typrest.HTTPError(err)
	}
	ec.Response().Header().Set(echo.HeaderLocation, fmt.Sprintf("/songs/%d", newSong.ID))
	return ec.JSON(http.StatusCreated, newSong)
}

// Find songs
func (c *SongCntrl) Find(ec echo.Context) (err error) {
	var req service.FindReq
	if err := ec.Bind(&req); err != nil {
		return err
	}
	ctx := ec.Request().Context()
	songs, err := c.SongSvc.Find(ctx, &req)
	if err != nil {
		return typrest.HTTPError(err)
	}
	return ec.JSON(http.StatusOK, songs)
}

// FindOne book
func (c *SongCntrl) FindOne(ec echo.Context) error {
	ctx := ec.Request().Context()
	id := ec.Param("id")
	book, err := c.SongSvc.FindOne(ctx, id)
	if err != nil {
		return typrest.HTTPError(err)
	}
	return ec.JSON(http.StatusOK, book)
}

// Delete book
func (c *SongCntrl) Delete(ec echo.Context) (err error) {
	if err = c.SongSvc.Delete(
		ec.Request().Context(),
		ec.Param("id"),
	); err != nil {
		return typrest.HTTPError(err)
	}
	return ec.NoContent(http.StatusNoContent)
}

// Update book
func (c *SongCntrl) Update(ec echo.Context) (err error) {
	var book mysqldb.Song
	if err = ec.Bind(&book); err != nil {
		return err
	}
	ctx := ec.Request().Context()
	paramID := ec.Param("id")
	updatedSong, err := c.SongSvc.Update(ctx, paramID, &book)
	if err != nil {
		return typrest.HTTPError(err)
	}
	return ec.JSON(http.StatusOK, updatedSong)
}

// Patch book
func (c *SongCntrl) Patch(ec echo.Context) (err error) {
	var book mysqldb.Song
	if err = ec.Bind(&book); err != nil {
		return err
	}
	ctx := ec.Request().Context()
	paramID := ec.Param("id")
	patchedSong, err := c.SongSvc.Patch(ctx, paramID, &book)
	if err != nil {
		return typrest.HTTPError(err)
	}
	return ec.JSON(http.StatusOK, patchedSong)
}
