package controller

import "github.com/labstack/echo"

// CRUD handle common create read update delete
type CRUD interface {
	// FIXME: separate between controller and CRUD function to different interface
	RegisterTo(entity string, e *echo.Echo)

	Create(c echo.Context) error
	Get(c echo.Context) error
	List(c echo.Context) error
	Update(c echo.Context) error
	Delete(c echo.Context) error

	BeforeActionFunc(next echo.HandlerFunc) echo.HandlerFunc
}
