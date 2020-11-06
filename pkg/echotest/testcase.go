package echotest

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/require"
)

type (
	// TestCase for echo
	TestCase struct {
		Request          Request
		ExpectedResponse Response
		ExpectedError    string
	}
	// Response that expected
	Response struct {
		Code   int
		Header http.Header
		Body   string
	}
	// Request for testcase
	Request struct {
		Method    string
		Target    string
		Body      string
		Header    http.Header
		URLParams map[string]string
	}
)

//
// TestCase
//

// Execute test case against handle function
func (tt *TestCase) Execute(t *testing.T, fn echo.HandlerFunc) {
	req := tt.Request.Request()
	rec, err := Do(fn, req, tt.Request.URLParams)
	if tt.ExpectedError != "" {
		require.EqualError(t, err, tt.ExpectedError)
	} else {
		require.NoError(t, err)
		EqualResp(t, tt.ExpectedResponse, rec)
	}
}

//
// Request
//

// Request return http request
func (r *Request) Request() *http.Request {
	req := httptest.NewRequest(r.Method, r.Target, strings.NewReader(r.Body))
	req.Header = r.Header
	return req
}
