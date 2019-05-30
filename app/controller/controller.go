package controller

import "github.com/labstack/echo"

// CRUDController handle common create read update delete
type CRUDController interface {
	Create(c echo.Context) error
	Get(c echo.Context) error
	List(c echo.Context) error
	Update(c echo.Context) error
	Delete(c echo.Context) error
}
