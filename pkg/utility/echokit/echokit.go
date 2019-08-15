package echokit

import (
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/labstack/echo"
)

func doRequest(handler echo.HandlerFunc, req *http.Request, urlParams map[string]string) (rec *httptest.ResponseRecorder, err error) {
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

// DoGET return echo.Context and httptest.ResponseRecorder for GET Request
func DoGET(handler echo.HandlerFunc, url string, urlParams map[string]string) (rec *httptest.ResponseRecorder, err error) {
	req := httptest.NewRequest(http.MethodGet, url, nil)
	return doRequest(handler, req, urlParams)
}

// DoPOST return echo.Context and httptest.ResponseRecorder for POST Request
func DoPOST(handler echo.HandlerFunc, url string, json string) (rec *httptest.ResponseRecorder, err error) {
	req := httptest.NewRequest(http.MethodPost, url, strings.NewReader(json))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	return doRequest(handler, req, nil)
}

// DoPUT return echo.Context and httptest.ResponseRecorder for POST Request
func DoPUT(handler echo.HandlerFunc, url string, json string) (rec *httptest.ResponseRecorder, err error) {
	req := httptest.NewRequest(http.MethodPut, url, strings.NewReader(json))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	return doRequest(handler, req, nil)
}

// DoDELETE return echo.Context and httptest.ResponseRecorder for DELETE Request
func DoDELETE(handler echo.HandlerFunc, url string, urlParams map[string]string) (rec *httptest.ResponseRecorder, err error) {
	req := httptest.NewRequest(http.MethodDelete, url, nil)
	return doRequest(handler, req, urlParams)
}
