package controller_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
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

func TestBookController_Get(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// prepare mock book repository
	bookR := mock.NewMockBookRepository(ctrl)
	bookR.EXPECT().Get(int64(1)).Return(&repository.Book{ID: 1, Title: "title1", Author: "author1"}, nil)
	bookR.EXPECT().Get(int64(2)).Return(nil, fmt.Errorf("some-get-error"))

	e := echo.New()
	defer e.Close()
	bookController := controller.NewBookController(bookR)
	bookController.RegisterTo("book", e)

	runTestCase(t, e, []ControllerTestCase{
		{
			"Invalid ID",
			RequestTestCase{http.MethodGet, "/book/abc", ""},
			ResponseTestCase{http.StatusBadRequest, "{\"message\":\"Invalid ID\"}\n"},
		},
		{
			"Get success",
			RequestTestCase{http.MethodGet, "/book/1", ""},
			ResponseTestCase{http.StatusOK, "{\"id\":1,\"title\":\"title1\",\"author\":\"author1\"}\n"},
		},
		{
			"Get error",
			RequestTestCase{http.MethodGet, "/book/2", ""},
			ResponseTestCase{http.StatusInternalServerError, "{\"message\":\"Internal Server Error\"}\n"},
		},
	})
}

func TestBookController_List(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// prepare mock book repository
	bookR := mock.NewMockBookRepository(ctrl)
	bookR.EXPECT().List().Return(
		[]*repository.Book{
			&repository.Book{ID: 1, Title: "title1", Author: "author1"},
			&repository.Book{ID: 2, Title: "title2", Author: "author2"},
		}, nil)
	bookR.EXPECT().List().Return(nil, fmt.Errorf("some-list-error"))

	e := echo.New()
	defer e.Close()
	bookController := controller.NewBookController(bookR)
	bookController.RegisterTo("book", e)

	runTestCase(t, e, []ControllerTestCase{
		{
			"List success",
			RequestTestCase{http.MethodGet, "/book", ""},
			ResponseTestCase{http.StatusOK, "[{\"id\":1,\"title\":\"title1\",\"author\":\"author1\"},{\"id\":2,\"title\":\"title2\",\"author\":\"author2\"}]\n"},
		},
		{
			"List error",
			RequestTestCase{http.MethodGet, "/book", ""},
			ResponseTestCase{http.StatusInternalServerError, "{\"message\":\"Internal Server Error\"}\n"},
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
	defer e.Close()
	bookController := controller.NewBookController(bookR)
	bookController.RegisterTo("book", e)

	runTestCase(t, e, []ControllerTestCase{
		{
			"Invalid message body",
			RequestTestCase{http.MethodPost, "/book", "{}"},
			ResponseTestCase{http.StatusBadRequest, "{\"message\":\"Invalid Message\"}\n"},
		},
		{
			"Invalid json format",
			RequestTestCase{http.MethodPost, "/book", "invalid-json"},
			ResponseTestCase{http.StatusBadRequest, "{\"message\":\"Syntax error: offset=1, error=invalid character 'i' looking for beginning of value\"}\n"},
		},
		{
			"Insert error",
			RequestTestCase{http.MethodPost, "/book", `{"author":"some-author", "title":"some-title"}`},
			ResponseTestCase{http.StatusInternalServerError, "{\"message\":\"Internal Server Error\"}\n"},
		},
		{
			"Insert Success",
			RequestTestCase{http.MethodPost, "/book", `{"author":"some-author", "title":"some-title"}`},
			ResponseTestCase{http.StatusCreated, "{\"message\":\"Success insert new record #99\"}\n"},
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
	defer e.Close()
	bookController := controller.NewBookController(bookR)
	bookController.RegisterTo("book", e)

	runTestCase(t, e, []ControllerTestCase{
		{
			"Invalid ID in url parameter",
			RequestTestCase{http.MethodDelete, "/book/abc", ``},
			ResponseTestCase{http.StatusBadRequest, "{\"message\":\"Invalid ID\"}\n"},
		},
		{
			"Valid ID",
			RequestTestCase{http.MethodDelete, "/book/1", ``},
			ResponseTestCase{http.StatusOK, "{\"message\":\"Delete #1 done\"}\n"},
		},
		{
			"ID not found",
			RequestTestCase{http.MethodDelete, "/book/2", ``},
			ResponseTestCase{http.StatusInternalServerError, "{\"message\":\"Internal Server Error\"}\n"},
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
	defer e.Close()
	bookController := controller.NewBookController(bookR)
	bookController.RegisterTo("book", e)

	runTestCase(t, e, []ControllerTestCase{
		{
			"Invalid json format",
			RequestTestCase{http.MethodPut, "/book", "invalid-json"},
			ResponseTestCase{http.StatusBadRequest, "{\"message\":\"Syntax error: offset=1, error=invalid character 'i' looking for beginning of value\"}\n"},
		},
		{
			"Invalid message body",
			RequestTestCase{http.MethodPut, "/book", `{}`},
			ResponseTestCase{http.StatusBadRequest, "{\"message\":\"Invalid ID\"}\n"},
		},
		{
			"Invalid message body",
			RequestTestCase{http.MethodPut, "/book", `{"id":1}`},
			ResponseTestCase{http.StatusBadRequest, "{\"message\":\"Invalid Message\"}\n"},
		},
		{
			"Update error",
			RequestTestCase{http.MethodPut, "/book", `{"id":1, "title":"some-title", "author": "some-author"}`},
			ResponseTestCase{http.StatusInternalServerError, "{\"message\":\"Internal Server Error\"}\n"},
		},
		{
			"Update success",
			RequestTestCase{http.MethodPut, "/book", `{"id":1, "title":"some-title", "author": "some-author"}`},
			ResponseTestCase{http.StatusOK, "{\"message\":\"Update success\"}\n"},
		},
	})
}
