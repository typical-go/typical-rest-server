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
		Header map[string]string
		Body   string
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
	req := httptest.NewRequest(tt.Request.Method, tt.Request.Target, strings.NewReader(tt.Request.Body))
	for key, value := range tt.Request.Header {
		req.Header.Set(key, value)
	}

	rec, err := Do(fn, req, tt.Request.URLParams)
	if tt.ExpectedError != "" {
		require.EqualError(t, err, tt.ExpectedError)
	} else {
		require.NoError(t, err)
		require.Equal(t, tt.ExpectedResponse.Code, rec.Code)
		require.Equal(t, tt.ExpectedResponse.Body, rec.Body.String())
		resHeader := rec.Result().Header
		for key, value := range tt.ExpectedResponse.Header {
			require.Equal(t, value, resHeader.Get(key))
		}
	}
}

// Do request against the handler
func Do(handler echo.HandlerFunc, req *http.Request, urlParams map[string]string) (rec *httptest.ResponseRecorder, err error) {
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
	return Do(handler, req, urlParams)
}

// DoPOST return echo.Context and httptest.ResponseRecorder for POST Request (deprecated)
func DoPOST(handler echo.HandlerFunc, url string, json string, urlParams map[string]string) (rec *httptest.ResponseRecorder, err error) {
	req := httptest.NewRequest(http.MethodPost, url, strings.NewReader(json))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	return Do(handler, req, urlParams)
}

// DoPUT return echo.Context and httptest.ResponseRecorder for PUT Request (deprecated)
func DoPUT(handler echo.HandlerFunc, url string, json string, urlParams map[string]string) (rec *httptest.ResponseRecorder, err error) {
	req := httptest.NewRequest(http.MethodPut, url, strings.NewReader(json))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	return Do(handler, req, urlParams)
}

// DoPATCH return echo.Context and httptest.ResponseRecorder for PUT Request (deprecated)
func DoPATCH(handler echo.HandlerFunc, url string, json string, urlParams map[string]string) (rec *httptest.ResponseRecorder, err error) {
	req := httptest.NewRequest(http.MethodPatch, url, strings.NewReader(json))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	return Do(handler, req, urlParams)
}

// DoDELETE return echo.Context and httptest.ResponseRecorder for DELETE Request (deprecated)
func DoDELETE(handler echo.HandlerFunc, url string, urlParams map[string]string) (rec *httptest.ResponseRecorder, err error) {
	req := httptest.NewRequest(http.MethodDelete, url, nil)
	return Do(handler, req, urlParams)
}
