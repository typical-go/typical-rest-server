package echotest

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/require"
)

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

// EqualResp expect response
func EqualResp(t *testing.T, expected Response, rec *httptest.ResponseRecorder) {
	require.Equal(t, expected.Code, rec.Code)
	require.Equal(t, expected.Body, rec.Body.String())
	require.Equal(t, expected.Header, rec.HeaderMap)
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
