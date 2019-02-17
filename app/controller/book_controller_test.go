package controller_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/imantung/typical-go-server/app/controller"
	"github.com/imantung/typical-go-server/app/repository"
	"github.com/labstack/echo"
	"github.com/stretchr/testify/require"
)

func TestBookController(t *testing.T) {

	testcases := []struct {
		repository   repository.BookRepository
		method       string
		target       string
		requestBody  io.Reader
		statusCode   int
		responseBody string
	}{
		{nil, http.MethodGet, "/book/1", nil, http.StatusInternalServerError, ""},
	}

	for i, tt := range testcases {
		req := httptest.NewRequest(tt.method, tt.target, tt.requestBody)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rr := httptest.NewRecorder()

		bookController := controller.NewBookController(tt.repository)

		e := echo.New()
		bookController.RegisterTo("book", e)
		http.HandlerFunc(e.ServeHTTP).ServeHTTP(rr, req)

		require.Equal(t, tt.statusCode, rr.Code, "Failed at test case %d", i)
	}
}
