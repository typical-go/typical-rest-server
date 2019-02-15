package controller

import "github.com/labstack/echo"

// CRUD handle common create read update delete
type CRUD interface {
	Create(c echo.Context) error
	Get(c echo.Context) error
	List(c echo.Context) error
	Update(c echo.Context) error
	Delete(c echo.Context) error
}
