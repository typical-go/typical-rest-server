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
	e.GET("{{.Table}}", c.Find)
	e.POST("{{.Table}}", c.Create)
	e.GET("{{.Table}}/:id", c.FindOne)
	e.PUT("{{.Table}}", c.Update)
	e.DELETE("{{.Table}}/:id", c.Delete)
}

// Create {{.Name}}
func (c *{{.Type}}Cntrl) Create(ec echo.Context) (err error) {
	var (
		{{.Name}} repository.{{.Type}}
		lastInsertID int64
		ctx = ec.Request().Context()
	)
	if err = ec.Bind(&{{.Name}}); err != nil {
		return err
	}
	if err = validator.New().Struct({{.Name}}); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	if lastInsertID, err = c.{{.Type}}Service.Create(ctx, {{.Name}}); err != nil {
		return echo.NewHTTPError(http.StatusUnprocessableEntity, err.Error())
	}
	return ec.JSON(http.StatusCreated, GeneralResponse{
		Message: fmt.Sprintf("Success insert new {{.Name}} #%d", lastInsertID),
	})
}

// Find of {{.Name}}
func (c *{{.Type}}Cntrl) Find(ec echo.Context) (err error) {
	var (
		{{.Name}}s []*repository.{{.Type}}
		ctx = ec.Request().Context()
	)
	if {{.Name}}s, err = c.{{.Type}}Service.Find(ctx); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return ec.JSON(http.StatusOK, {{.Name}}s)
}

// FindOne {{.Name}}
func (c *{{.Type}}Cntrl) FindOne(ec echo.Context) (err error) {
	var (
		id int64
		{{.Name}} *repository.{{.Type}}
		ctx = ec.Request().Context()
	)
	if id, err = strconv.ParseInt(ec.Param("id"), 10, 64); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid ID")
	}
	if {{.Name}}, err = c.{{.Type}}Service.FindOne(ctx, id); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	if {{.Name}} == nil {
		return echo.NewHTTPError(http.StatusNotFound, fmt.Sprintf("{{.Type}}#%d not found", id))
	}
	return ec.JSON(http.StatusOK, {{.Name}})
}

// Delete {{.Name}}
func (c *{{.Type}}Cntrl) Delete(ec echo.Context) (err error) {
	var (
		id int64
		ctx = ec.Request().Context()
	)
	if id, err = strconv.ParseInt(ec.Param("id"), 10, 64); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid ID")
	}
	if err = c.{{.Type}}Service.Delete(ctx, id); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return ec.JSON(http.StatusOK, GeneralResponse{
		Message: fmt.Sprintf("Success delete {{.Name}} #%d", id),
	})
}

// Update {{.Name}}
func (c *{{.Type}}Cntrl) Update(ec echo.Context) (err error) {
	var (
		{{.Name}} repository.{{.Type}}
		ctx = ec.Request().Context()
	)
	if err = ec.Bind(&{{.Name}}); err != nil {
		return err
	}
	if {{.Name}}.ID <= 0 {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid ID")
	}
	if err = validator.New().Struct({{.Name}}); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	if err = c.{{.Type}}Service.Update(ctx, {{.Name}}); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return ec.JSON(http.StatusOK, GeneralResponse{
		Message: fmt.Sprintf("Success update {{.Name}} #%d", {{.Name}}.ID),
	})
}
`
