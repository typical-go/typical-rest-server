package controller_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/imantung/typical-go-server/app/controller"
	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"
)

func TestBookController(t *testing.T) {

	bookController := controller.NewBookController()

	testcases := []struct {
		handler      func(e echo.Context) error
		method       string
		target       string
		requestBody  io.Reader
		statusCode   int
		responseBody string
	}{
		{bookController.List, http.MethodGet, "/book", nil, http.StatusServiceUnavailable, ""},
		{bookController.Create, http.MethodPost, "/book", nil, http.StatusServiceUnavailable, ""},
		{bookController.Get, http.MethodGet, "/book/1", nil, http.StatusServiceUnavailable, ""},
		{bookController.Delete, http.MethodDelete, "/book/1", nil, http.StatusServiceUnavailable, ""},
		{bookController.Update, http.MethodPut, "/book", nil, http.StatusServiceUnavailable, ""},
	}

	e := echo.New()

	for _, tt := range testcases {
		req := httptest.NewRequest(tt.method, tt.target, tt.requestBody)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		// Assertions
		if assert.NoError(t, tt.handler(c)) {
			assert.Equal(t, tt.statusCode, rec.Code)
		}
	}
}
