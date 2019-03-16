package controller_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/imantung/typical-go-server/app/controller"
	"github.com/imantung/typical-go-server/app/helper/timekit"
	"github.com/imantung/typical-go-server/app/repository"
	"github.com/imantung/typical-go-server/test/mock"
	"github.com/labstack/echo"
	"github.com/stretchr/testify/require"
)

// dummy data
var (
	book1 = &repository.Book{ID: 1, Title: "title1", Author: "author1", CreatedAt: timekit.UTC("2019-02-20T10:00:00-05:00")}
	book2 = &repository.Book{ID: 2, Title: "title2", Author: "author2", CreatedAt: timekit.UTC("2019-02-20T10:00:01-05:00")}
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

func TestBookController_Get(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// prepare mock book repository
	bookR := mock.NewMockBookRepository(ctrl)
	bookR.EXPECT().Get(int64(1)).Return(book1, nil)
	bookR.EXPECT().Get(int64(2)).Return(nil, fmt.Errorf("some-get-error"))

	e := echo.New()

	bookController := controller.NewBookController(bookR)
	bookController.RegisterTo("book", e)

	runTestCase(t, e, []ControllerTestCase{
		{
			"Invalid ID",
			http.MethodGet, "/book/abc", "",
			http.StatusBadRequest, "{\"message\":\"Invalid ID\"}\n",
		},
		{
			"Get success",
			http.MethodGet, "/book/1", "",
			http.StatusOK, "{\"id\":1,\"title\":\"title1\",\"author\":\"author1\"}\n",
		},
		{
			"Get error",
			http.MethodGet, "/book/2", "",
			http.StatusInternalServerError, "{\"message\":\"Internal Server Error\"}\n",
		},
	})
}

func TestBookController_List(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// prepare mock book repository
	bookR := mock.NewMockBookRepository(ctrl)
	bookR.EXPECT().List().Return([]*repository.Book{book1, book2}, nil)
	bookR.EXPECT().List().Return(nil, fmt.Errorf("some-list-error"))

	e := echo.New()

	bookController := controller.NewBookController(bookR)
	bookController.RegisterTo("book", e)

	runTestCase(t, e, []ControllerTestCase{
		{
			"List success",
			http.MethodGet, "/book", "",
			http.StatusOK, "[{\"id\":1,\"title\":\"title1\",\"author\":\"author1\"},{\"id\":2,\"title\":\"title2\",\"author\":\"author2\"}]\n",
		},
		{
			"List error",
			http.MethodGet, "/book", "",
			http.StatusInternalServerError, "{\"message\":\"Internal Server Error\"}\n",
		},
	})
}

func TestBookController_Insert(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// prepare mock book repository
	bookR := mock.NewMockBookRepository(ctrl)
	bookR.EXPECT().Insert(gomock.Any()).Return(int64(0), fmt.Errorf("some-insert-error"))
	bookR.EXPECT().Insert(gomock.Any()).Return(int64(99), nil)

	e := echo.New()

	bookController := controller.NewBookController(bookR)
	bookController.RegisterTo("book", e)

	runTestCase(t, e, []ControllerTestCase{
		{
			"Invalid message body",
			http.MethodPost, "/book", "{}",
			http.StatusBadRequest, "{\"message\":\"Invalid Message\"}\n",
		},
		{
			"Invalid json format",
			http.MethodPost, "/book", "invalid-json",
			http.StatusBadRequest, "{\"message\":\"Syntax error: offset=1, error=invalid character 'i' looking for beginning of value\"}\n",
		},
		{
			"Insert error",
			http.MethodPost, "/book", `{"author":"some-author", "title":"some-title"}`,
			http.StatusInternalServerError, "{\"message\":\"Internal Server Error\"}\n",
		},
		{
			"Insert Success",
			http.MethodPost, "/book", `{"author":"some-author", "title":"some-title"}`,
			http.StatusCreated, "{\"message\":\"Success insert new record #99\"}\n",
		},
	})
}

func TestBookController_Delete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// prepare mock book repository
	bookR := mock.NewMockBookRepository(ctrl)
	bookR.EXPECT().Delete(int64(1)).Return(nil)
	bookR.EXPECT().Delete(int64(2)).Return(fmt.Errorf("some-delete-error"))

	e := echo.New()

	bookController := controller.NewBookController(bookR)
	bookController.RegisterTo("book", e)

	runTestCase(t, e, []ControllerTestCase{
		{
			"Invalid ID in url parameter",
			http.MethodDelete, "/book/abc", ``,
			http.StatusBadRequest, "{\"message\":\"Invalid ID\"}\n",
		},
		{
			"Valid ID",
			http.MethodDelete, "/book/1", ``,
			http.StatusOK, "{\"message\":\"Delete #1 done\"}\n",
		},
		{
			"ID not found",
			http.MethodDelete, "/book/2", ``,
			http.StatusInternalServerError, "{\"message\":\"Internal Server Error\"}\n",
		},
	})
}

func TestBookController_Update(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// prepare mock book repository
	bookR := mock.NewMockBookRepository(ctrl)
	bookR.EXPECT().Update(gomock.Any()).Return(fmt.Errorf("some-update-error"))
	bookR.EXPECT().Update(gomock.Any()).Return(nil)

	e := echo.New()

	bookController := controller.NewBookController(bookR)
	bookController.RegisterTo("book", e)

	runTestCase(t, e, []ControllerTestCase{
		{
			"Invalid json format",
			http.MethodPut, "/book", "invalid-json",
			http.StatusBadRequest, "{\"message\":\"Syntax error: offset=1, error=invalid character 'i' looking for beginning of value\"}\n",
		},
		{
			"Invalid message body",
			http.MethodPut, "/book", `{}`,
			http.StatusBadRequest, "{\"message\":\"Invalid ID\"}\n",
		},
		{
			"Invalid message body",
			http.MethodPut, "/book", `{"id":1}`,
			http.StatusBadRequest, "{\"message\":\"Invalid Message\"}\n",
		},
		{
			"Update error",
			http.MethodPut, "/book", `{"id":1, "title":"some-title", "author": "some-author"}`,
			http.StatusInternalServerError, "{\"message\":\"Internal Server Error\"}\n",
		},
		{
			"Update success",
			http.MethodPut, "/book", `{"id":1, "title":"some-title", "author": "some-author"}`,
			http.StatusOK, "{\"message\":\"Update success\"}\n",
		},
	})
}
