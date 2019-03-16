package controller_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo"
	"github.com/stretchr/testify/require"
)

type RequestTestCase struct {
	method      string
	target      string
	requestBody string
}

type ResponseTestCase struct {
	statusCode   int
	responseBody string
}

type ControllerTestCase struct {
	description string
	RequestTestCase
	ResponseTestCase
}

func runTestCase(t *testing.T, handler http.Handler, testcases []ControllerTestCase) {
	for _, tt := range testcases {
		req := httptest.NewRequest(tt.method, tt.target, strings.NewReader(tt.requestBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rr := httptest.NewRecorder()
		http.HandlerFunc(handler.ServeHTTP).ServeHTTP(rr, req)

		require.Equal(t, tt.statusCode, rr.Code, "%s: Expect %d but got %d",
			tt.description, tt.statusCode, rr.Code)
		require.Equal(t, tt.responseBody, rr.Body.String(), "%s: Expect %s but got %s",
			tt.description, tt.responseBody, rr.Body.String())
	}
}
