package controller_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/imantung/go-helper/timekit"
	"github.com/imantung/typical-go-server/app/controller"
	"github.com/imantung/typical-go-server/app/repository"
	"github.com/imantung/typical-go-server/test/mock"
	"github.com/labstack/echo"
	"github.com/stretchr/testify/require"
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

	bookR := mock.NewMockBookRepository(ctrl)
	bookR.EXPECT().Get(1).Return(
		repository.Book{ID: 1, Title: "title1", Author: "author1", CreatedAt: timekit.UTC("2019-02-20T10:00:00-05:00")},
		nil,
	)
	bookR.EXPECT().Get(2).Return(
		repository.Book{},
		fmt.Errorf("some-error"),
	)

	bookController := controller.NewBookController(bookR)
	e := echo.New()
	bookController.RegisterTo("book", e)

	testcases := []struct {
		method       string
		target       string
		requestBody  string
		statusCode   int
		responseBody string
	}{
		{
			http.MethodGet, "/book/1", "",
			http.StatusOK, `{"id":1,"title":"title1","author":"author1","created_at":"2019-02-20T10:00:00-05:00"}` + "\n",
		},
		{
			http.MethodGet, "/book/2", "",
			http.StatusInternalServerError, `{"message":"Internal Server Error"}` + "\n",
		},
	}

	for i, tt := range testcases {
		rr := httptest.NewRecorder()
		http.HandlerFunc(e.ServeHTTP).ServeHTTP(rr,
			httptest.NewRequest(tt.method, tt.target, strings.NewReader(tt.requestBody)))

		require.Equal(t, tt.statusCode, rr.Code, "%d: Expect %d but got %d",
			i, tt.statusCode, rr.Code)
		require.Equal(t, tt.responseBody, rr.Body.String(), "%d: Expect %s but got %s",
			i, tt.responseBody, rr.Body.String())
	}
}
