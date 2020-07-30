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
		Request

		ExpectedCode   int
		ExpectedHeader map[string]string
		ExpectedBody   string
		ExpectedErr    string
	}

	// Request for testcase
	Request struct {
		Method    string
		Target    string
		Body      string
		Header    map[string]string
		URLParams map[string]string
	}
)

// Execute test case against handle function
func (tt *TestCase) Execute(t *testing.T, fn echo.HandlerFunc) {
	req := httptest.NewRequest(tt.Method, tt.Target, strings.NewReader(tt.Body))
	for key, value := range tt.Header {
		req.Header.Set(key, value)
	}

	rec, err := execute(fn, req, tt.URLParams)
	if tt.ExpectedErr != "" {
		require.EqualError(t, err, tt.ExpectedErr)

		return
	}

	require.NoError(t, err)
	require.Equal(t, tt.ExpectedCode, rec.Code)
	require.Equal(t, tt.ExpectedBody, rec.Body.String())
	resHeader := rec.Result().Header
	for key, value := range tt.ExpectedHeader {
		require.Equal(t, value, resHeader.Get(key))
	}

}

func execute(handler echo.HandlerFunc, req *http.Request, urlParams map[string]string) (rec *httptest.ResponseRecorder, err error) {
	rec = httptest.NewRecorder()
	e := echo.New()
	ctx := e.NewContext(req, rec)
	if urlParams != nil {
		var keys []string
		var values []string
		for key, value := range urlParams {
			keys = append(keys, key)
			values = append(values, value)
		}
		ctx.SetParamNames(keys...)
		ctx.SetParamValues(values...)
	}
	err = handler(ctx)
	return
}

// HeaderForJSON to generate header for json-based API
func HeaderForJSON() map[string]string {
	return map[string]string{
		echo.HeaderContentType: echo.MIMEApplicationJSON,
	}
}

// DoGET return recorder and error for GET API (deprecated)
func DoGET(handler echo.HandlerFunc, url string, urlParams map[string]string) (rec *httptest.ResponseRecorder, err error) {
	req := httptest.NewRequest(http.MethodGet, url, nil)
	return execute(handler, req, urlParams)
}

// DoPOST return echo.Context and httptest.ResponseRecorder for POST Request (deprecated)
func DoPOST(handler echo.HandlerFunc, url string, json string, urlParams map[string]string) (rec *httptest.ResponseRecorder, err error) {
	req := httptest.NewRequest(http.MethodPost, url, strings.NewReader(json))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	return execute(handler, req, urlParams)
}

// DoPUT return echo.Context and httptest.ResponseRecorder for PUT Request (deprecated)
func DoPUT(handler echo.HandlerFunc, url string, json string, urlParams map[string]string) (rec *httptest.ResponseRecorder, err error) {
	req := httptest.NewRequest(http.MethodPut, url, strings.NewReader(json))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	return execute(handler, req, urlParams)
}

// DoDELETE return echo.Context and httptest.ResponseRecorder for DELETE Request (deprecated)
func DoDELETE(handler echo.HandlerFunc, url string, urlParams map[string]string) (rec *httptest.ResponseRecorder, err error) {
	req := httptest.NewRequest(http.MethodDelete, url, nil)
	return execute(handler, req, urlParams)
}
