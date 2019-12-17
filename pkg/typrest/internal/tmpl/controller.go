package tmpl

// Controller template
const Controller = `package controller

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/labstack/echo"

	"{{.ProjectPackage}}/app/repository"
	"{{.ProjectPackage}}/app/service"
	"go.uber.org/dig"
	"gopkg.in/go-playground/validator.v9"
)

// {{.Type}}Cntrl is controller to {{.Name}} entity
type {{.Type}}Cntrl struct {
	dig.In
	service.{{.Type}}Service
}

// Route to define API Route
func (c *{{.Type}}Cntrl) Route(e *echo.Echo) {
	e.GET("{{.Table}}", c.List)
	e.POST("{{.Table}}", c.Create)
	e.GET("{{.Table}}/:id", c.Get)
	e.PUT("{{.Table}}", c.Update)
	e.DELETE("{{.Table}}/:id", c.Delete)
}

// Create {{.Name}}
func (c *{{.Type}}Cntrl) Create(ctx echo.Context) (err error) {
	var {{.Name}} repository.{{.Type}}
	var lastInsertID int64
	ctx0 := ctx.Request().Context()
	if err = ctx.Bind(&{{.Name}}); err != nil {
		return err
	}
	if err = validator.New().Struct({{.Name}}); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	if lastInsertID, err = c.{{.Type}}Service.Insert(ctx0, {{.Name}}); err != nil {
		return echo.NewHTTPError(http.StatusUnprocessableEntity, err.Error())
	}
	return ctx.JSON(http.StatusCreated, GeneralResponse{
		Message: fmt.Sprintf("Success insert new {{.Name}} #%d", lastInsertID),
	})
}

// List of {{.Name}}
func (c *{{.Type}}Cntrl) List(ctx echo.Context) (err error) {
	var {{.Name}}s []*repository.{{.Type}}
	ctx0 := ctx.Request().Context()
	if {{.Name}}s, err = c.{{.Type}}Service.List(ctx0); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return ctx.JSON(http.StatusOK, {{.Name}}s)
}

// Get {{.Name}}
func (c *{{.Type}}Cntrl) Get(ctx echo.Context) (err error) {
	var id int64
	var {{.Name}} *repository.{{.Type}}
	ctx0 := ctx.Request().Context()
	if id, err = strconv.ParseInt(ctx.Param("id"), 10, 64); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid ID")
	}
	if {{.Name}}, err = c.{{.Type}}Service.Find(ctx0, id); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	if {{.Name}} == nil {
		return echo.NewHTTPError(http.StatusNotFound, fmt.Sprintf("{{.Type}}#%d not found", id))
	}
	return ctx.JSON(http.StatusOK, {{.Name}})
}

// Delete {{.Name}}
func (c *{{.Type}}Cntrl) Delete(ctx echo.Context) (err error) {
	var id int64
	ctx0 := ctx.Request().Context()
	if id, err = strconv.ParseInt(ctx.Param("id"), 10, 64); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid ID")
	}
	if err = c.{{.Type}}Service.Delete(ctx0, id); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return ctx.JSON(http.StatusOK, GeneralResponse{
		Message: fmt.Sprintf("Success delete {{.Name}} #%d", id),
	})
}

// Update {{.Name}}
func (c *{{.Type}}Cntrl) Update(ctx echo.Context) (err error) {
	var {{.Name}} repository.{{.Type}}
	ctx0 := ctx.Request().Context()
	if err = ctx.Bind(&{{.Name}}); err != nil {
		return err
	}
	if {{.Name}}.ID <= 0 {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid ID")
	}
	if err = validator.New().Struct({{.Name}}); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	if err = c.{{.Type}}Service.Update(ctx0, {{.Name}}); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return ctx.JSON(http.StatusOK, GeneralResponse{
		Message: fmt.Sprintf("Success update {{.Name}} #%d", {{.Name}}.ID),
	})
}
`
