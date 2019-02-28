package controller_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/golang/mock/gomock"
	"github.com/imantung/go-helper/timekit"
	"github.com/imantung/typical-go-server/app/controller"
	"github.com/imantung/typical-go-server/app/repository"
	"github.com/imantung/typical-go-server/test/mock"
	"github.com/labstack/echo"
	"github.com/stretchr/testify/require"
)

// dummy data
var (
	book1 = repository.Book{ID: 1, Title: "title1", Author: "author1", CreatedAt: timekit.UTC("2019-02-20T10:00:00-05:00")}
	book2 = repository.Book{ID: 2, Title: "title2", Author: "author2", CreatedAt: timekit.UTC("2019-02-20T10:00:01-05:00")}
)

func TestBookController_NoRepository(t *testing.T) {
	e := echo.New()

	bookController := controller.NewBookController(nil)
	bookController.RegisterTo("book", e)

	rr := httptest.NewRecorder()
	http.HandlerFunc(e.ServeHTTP).ServeHTTP(rr,
		httptest.NewRequest(http.MethodGet, "/book/1", nil))

	require.Equal(t, http.StatusInternalServerError, rr.Code)
}

func TestBookController(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// prepare mock book repository
	bookR := mock.NewMockBookRepository(ctrl)
	bookR.EXPECT().Get(1).Return(book1, nil)
	bookR.EXPECT().Get(2).Return(repository.Book{}, fmt.Errorf("some-error"))
	bookR.EXPECT().List().Return([]repository.Book{book1, book2}, nil)
	bookR.EXPECT().List().Return(nil, fmt.Errorf("some-error"))
	bookR.EXPECT().Insert(gomock.Any()).Return(nil, fmt.Errorf("some-error"))
	bookR.EXPECT().Insert(gomock.Any()).Return(sqlmock.NewResult(99, 1), nil)

	bookController := controller.NewBookController(bookR)
	e := echo.New()
	bookController.RegisterTo("book", e)

	// create testcases
	testcases := []struct {
		method       string
		target       string
		requestBody  string
		statusCode   int
		responseBody string
	}{
		{
			http.MethodGet, "/book/1", "",
			http.StatusOK, "{\"id\":1,\"title\":\"title1\",\"author\":\"author1\",\"created_at\":\"2019-02-20T10:00:00-05:00\"}\n",
		},
		{
			http.MethodGet, "/book/2", "",
			http.StatusInternalServerError, "{\"message\":\"Internal Server Error\"}\n",
		},
		{
			http.MethodGet, "/book", "",
			http.StatusOK, "[{\"id\":1,\"title\":\"title1\",\"author\":\"author1\",\"created_at\":\"2019-02-20T10:00:00-05:00\"},{\"id\":2,\"title\":\"title2\",\"author\":\"author2\",\"created_at\":\"2019-02-20T10:00:01-05:00\"}]\n",
		},
		{
			http.MethodGet, "/book", "",
			http.StatusInternalServerError, "{\"message\":\"Internal Server Error\"}\n",
		},
		{
			http.MethodPost, "/book", "{}",
			http.StatusBadRequest, "{\"message\":\"Invalid Message\"}\n",
		},
		{
			http.MethodPost, "/book", "invalid-json",
			http.StatusBadRequest, "{\"message\":\"Syntax error: offset=1, error=invalid character 'i' looking for beginning of value\"}\n",
		},
		{
			http.MethodPost, "/book", `{"author":"some-author", "title":"some-title"}`,
			http.StatusInternalServerError, "{\"message\":\"Internal Server Error\"}\n",
		},
		{
			http.MethodPost, "/book", `{"author":"some-author", "title":"some-title"}`,
			http.StatusCreated, "{\"message\":\"Success insert new record #99\"}\n",
		},
	}

	// execute the test
	for i, tt := range testcases {
		req := httptest.NewRequest(tt.method, tt.target, strings.NewReader(tt.requestBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rr := httptest.NewRecorder()
		http.HandlerFunc(e.ServeHTTP).ServeHTTP(rr, req)

		require.Equal(t, tt.statusCode, rr.Code, "%d: Expect %d but got %d",
			i, tt.statusCode, rr.Code)
		require.Equal(t, tt.responseBody, rr.Body.String(), "%d: Expect %s but got %s",
			i, tt.responseBody, rr.Body.String())
	}
}
