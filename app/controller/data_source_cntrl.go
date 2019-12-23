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

// DataSourceCntrl is controller to data_source entity
type DataSourceCntrl struct {
	dig.In
	service.DataSourceService
}

// Route to define API Route
func (c *DataSourceCntrl) Route(e *echo.Echo) {
	e.GET("data_sources", c.List)
	e.POST("data_sources", c.Create)
	e.GET("data_sources/:id", c.Get)
	e.PUT("data_sources", c.Update)
	e.DELETE("data_sources/:id", c.Delete)
}

// Create data_source
func (c *DataSourceCntrl) Create(ctx echo.Context) (err error) {
	var data_source repository.DataSource
	var lastInsertID int64
	ctx0 := ctx.Request().Context()
	if err = ctx.Bind(&data_source); err != nil {
		return err
	}
	if err = validator.New().Struct(data_source); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	if lastInsertID, err = c.DataSourceService.Insert(ctx0, data_source); err != nil {
		return echo.NewHTTPError(http.StatusUnprocessableEntity, err.Error())
	}
	return ctx.JSON(http.StatusCreated, GeneralResponse{
		Message: fmt.Sprintf("Success insert new data_source #%d", lastInsertID),
	})
}

// List of data_source
func (c *DataSourceCntrl) List(ctx echo.Context) (err error) {
	var data_sources []*repository.DataSource
	ctx0 := ctx.Request().Context()
	if data_sources, err = c.DataSourceService.List(ctx0); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return ctx.JSON(http.StatusOK, data_sources)
}

// Get data_source
func (c *DataSourceCntrl) Get(ctx echo.Context) (err error) {
	var id int64
	var data_source *repository.DataSource
	ctx0 := ctx.Request().Context()
	if id, err = strconv.ParseInt(ctx.Param("id"), 10, 64); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid ID")
	}
	if data_source, err = c.DataSourceService.Find(ctx0, id); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	if data_source == nil {
		return echo.NewHTTPError(http.StatusNotFound, fmt.Sprintf("DataSource#%d not found", id))
	}
	return ctx.JSON(http.StatusOK, data_source)
}

// Delete data_source
func (c *DataSourceCntrl) Delete(ctx echo.Context) (err error) {
	var id int64
	ctx0 := ctx.Request().Context()
	if id, err = strconv.ParseInt(ctx.Param("id"), 10, 64); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid ID")
	}
	if err = c.DataSourceService.Delete(ctx0, id); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return ctx.JSON(http.StatusOK, GeneralResponse{
		Message: fmt.Sprintf("Success delete data_source #%d", id),
	})
}

// Update data_source
func (c *DataSourceCntrl) Update(ctx echo.Context) (err error) {
	var data_source repository.DataSource
	ctx0 := ctx.Request().Context()
	if err = ctx.Bind(&data_source); err != nil {
		return err
	}
	if data_source.ID <= 0 {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid ID")
	}
	if err = validator.New().Struct(data_source); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	if err = c.DataSourceService.Update(ctx0, data_source); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return ctx.JSON(http.StatusOK, GeneralResponse{
		Message: fmt.Sprintf("Success update data_source #%d", data_source.ID),
	})
}
