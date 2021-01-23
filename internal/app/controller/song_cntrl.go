package controller

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/typical-go/typical-rest-server/internal/app/service"
	"github.com/typical-go/typical-rest-server/internal/app/entity"
	"github.com/typical-go/typical-rest-server/pkg/cachekit"
	"github.com/typical-go/typical-rest-server/pkg/echokit"
	"go.uber.org/dig"
)

type (
	// SongCntrl is controller to book entity
	SongCntrl struct {
		dig.In
		Cache *cachekit.Store
		Svc   service.SongSvc
	}
)

var _ echokit.Router = (*SongCntrl)(nil)

// SetRoute to define API Route
func (c *SongCntrl) SetRoute(e echokit.Server) {
	e.GET("/songs", c.Find, c.Cache.Middleware)
	e.GET("/songs/:id", c.FindOne, c.Cache.Middleware)
	e.HEAD("/songs/:id", c.FindOne, c.Cache.Middleware)
	e.POST("/songs", c.Create)
	e.PUT("/songs/:id", c.Update)
	e.PATCH("/songs/:id", c.Patch)
	e.DELETE("/songs/:id", c.Delete)
}

// Create book
func (c *SongCntrl) Create(ec echo.Context) (err error) {
	var book entity.Song
	if err = ec.Bind(&book); err != nil {
		return err
	}
	ctx := ec.Request().Context()
	newSong, err := c.Svc.Create(ctx, &book)
	if err != nil {
		return echokit.HTTPError(err)
	}
	ec.Response().Header().Set(echo.HeaderLocation, fmt.Sprintf("/songs/%d", newSong.ID))
	return ec.JSON(http.StatusCreated, newSong)
}

// Find songs
func (c *SongCntrl) Find(ec echo.Context) (err error) {
	var req service.FindSongReq
	if err := ec.Bind(&req); err != nil {
		return err
	}
	ctx := ec.Request().Context()
	resp, err := c.Svc.Find(ctx, &req)
	if err != nil {
		return echokit.HTTPError(err)
	}
	ec.Response().Header().Add(echokit.HeaderTotalCount, resp.TotalCount)
	return ec.JSON(http.StatusOK, resp.Songs)
}

// FindOne book
func (c *SongCntrl) FindOne(ec echo.Context) error {
	ctx := ec.Request().Context()
	id := ec.Param("id")
	book, err := c.Svc.FindOne(ctx, id)
	if err != nil {
		return echokit.HTTPError(err)
	}
	return ec.JSON(http.StatusOK, book)
}

// Delete book
func (c *SongCntrl) Delete(ec echo.Context) (err error) {
	if err = c.Svc.Delete(
		ec.Request().Context(),
		ec.Param("id"),
	); err != nil {
		return echokit.HTTPError(err)
	}
	return ec.NoContent(http.StatusNoContent)
}

// Update book
func (c *SongCntrl) Update(ec echo.Context) (err error) {
	var book entity.Song
	if err = ec.Bind(&book); err != nil {
		return err
	}
	ctx := ec.Request().Context()
	paramID := ec.Param("id")
	updatedSong, err := c.Svc.Update(ctx, paramID, &book)
	if err != nil {
		return echokit.HTTPError(err)
	}
	return ec.JSON(http.StatusOK, updatedSong)
}

// Patch book
func (c *SongCntrl) Patch(ec echo.Context) (err error) {
	var book entity.Song
	if err = ec.Bind(&book); err != nil {
		return err
	}
	ctx := ec.Request().Context()
	paramID := ec.Param("id")
	patchedSong, err := c.Svc.Patch(ctx, paramID, &book)
	if err != nil {
		return echokit.HTTPError(err)
	}
	return ec.JSON(http.StatusOK, patchedSong)
}
