package typrest

import (
	"database/sql"
	"net/http"

	"github.com/labstack/echo/v4"
)

type (
	// Error interface
	Error interface {
		error
		HTTPError() *echo.HTTPError
	}
	errorImpl struct {
		code    int
		message string
	}
	// ValidErr validation error
	ValidErr struct {
		Message string
	}
)

// HTTPError convert error to *echo.HTTPError
func HTTPError(err error) *echo.HTTPError {
	if err == sql.ErrNoRows {
		return echo.NewHTTPError(http.StatusNotFound)
	}
	if restErr, ok := err.(Error); ok {
		return restErr.HTTPError()
	}
	if httpErr, ok := err.(*echo.HTTPError); ok {
		return httpErr
	}
	return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
}

//
// NewError
//

// NewError return error
func NewError(code int, message string) Error {
	return &errorImpl{code: code, message: message}
}

// NotImplementedErr return not implemented error
func NotImplementedErr() Error {
	return NewError(http.StatusNotImplemented, "not implemented")
}

// HTTPError return echo HTTP Error
func (v *errorImpl) HTTPError() *echo.HTTPError {
	return echo.NewHTTPError(v.code, v.message)
}

func (v *errorImpl) Error() string {
	return v.message
}

//
// ValidationError
//

var _ (Error) = (*ValidErr)(nil)

// NewValidErr create ValidationError
func NewValidErr(message string) *ValidErr {
	return &ValidErr{Message: message}
}

// HTTPError return echo HTTP Error
func (v *ValidErr) HTTPError() *echo.HTTPError {
	return echo.NewHTTPError(http.StatusUnprocessableEntity, v.Message)
}

func (v *ValidErr) Error() string {
	return v.Message
}
