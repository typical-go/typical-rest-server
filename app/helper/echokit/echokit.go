package echokit

import (
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/labstack/echo"
)

// Request return echo.Context and httptest.ResponseRecorder
func Request(req *http.Request) (ctx echo.Context, rec *httptest.ResponseRecorder) {
	rec = httptest.NewRecorder()

	e := echo.New()
	ctx = e.NewContext(req, rec)
	return
}

// RequestGET return echo.Context and httptest.ResponseRecorder for GET Request
func RequestGET(url string) (ctx echo.Context, rec *httptest.ResponseRecorder) {
	ctx, rec = Request(httptest.NewRequest(http.MethodGet, url, nil))
	return
}

// RequestGETWithParam return echo.Context and httptest.ResponseRecorder for GET Request with URL Param
func RequestGETWithParam(url string, urlParams map[string]string) (ctx echo.Context, rec *httptest.ResponseRecorder) {
	ctx, rec = RequestGET(url)
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
	return
}

// RequestPOST return echo.Context and httptest.ResponseRecorder for POST Request
func RequestPOST(url string, json string) (ctx echo.Context, rec *httptest.ResponseRecorder) {
	req := httptest.NewRequest(http.MethodPost, url, strings.NewReader(json))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	return Request(req)
}

// RequestPUT return echo.Context and httptest.ResponseRecorder for POST Request
func RequestPUT(url string, json string) (ctx echo.Context, rec *httptest.ResponseRecorder) {
	req := httptest.NewRequest(http.MethodPut, url, strings.NewReader(json))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	return Request(req)
}

// RequestDELETE return echo.Context and httptest.ResponseRecorder for DELETE Request
func RequestDELETE(url string) (ctx echo.Context, rec *httptest.ResponseRecorder) {
	ctx, rec = Request(httptest.NewRequest(http.MethodDelete, url, nil))
	return
}

// RequestDELETEWithParam return echo.Context and httptest.ResponseRecorder for DELETE Request with URL Param
func RequestDELETEWithParam(url string, urlParams map[string]string) (ctx echo.Context, rec *httptest.ResponseRecorder) {
	ctx, rec = RequestDELETE(url)
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
	return
}
