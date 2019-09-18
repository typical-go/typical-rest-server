package controller_test

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/typical-go/typical-rest-server/app/controller"
	"github.com/typical-go/typical-rest-server/app/repository"
	"github.com/typical-go/typical-rest-server/pkg/utility/echokit"

	"github.com/typical-go/typical-rest-server/mock"

	"github.com/stretchr/testify/require"
)

func TestBookController_Get(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	bookSvc := mock.NewMockBookService(ctrl)
	bookCntrl := controller.BookController{
		BookService: bookSvc,
	}
	t.Run("GIVEN invalid id", func(t *testing.T) {
		rr, err := echokit.DoGET(bookCntrl.Get, "/", map[string]string{
			"id": "invalid",
		})
		require.NoError(t, err)
		require.Equal(t, http.StatusBadRequest, rr.Code)
	})

	t.Run("GIVEN valid ID", func(t *testing.T) {
		bookSvc.EXPECT().Find(gomock.Any(), int64(1)).Return(&repository.Book{ID: 1, Title: "title1", Author: "author1"}, nil)

		rr, err := echokit.DoGET(bookCntrl.Get, "/", map[string]string{
			"id": "1",
		})
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, rr.Code)
		require.Equal(t, "{\"id\":1,\"title\":\"title1\",\"author\":\"author1\"}\n", rr.Body.String())
	})

	t.Run("When repository not found", func(t *testing.T) {
		bookSvc.EXPECT().Find(gomock.Any(), int64(3)).Return(nil, nil)

		rr, err := echokit.DoGET(bookCntrl.Get, "/", map[string]string{
			"id": "3",
		})
		require.NoError(t, err)
		require.Equal(t, http.StatusNotFound, rr.Code)
		require.Equal(t, "{\"message\":\"book #3 not found\"}\n", rr.Body.String())
	})

	t.Run("When return error", func(t *testing.T) {
		bookSvc.EXPECT().Find(gomock.Any(), int64(2)).Return(nil, fmt.Errorf("some-get-error"))

		_, err := echokit.DoGET(bookCntrl.Get, "/", map[string]string{
			"id": "2",
		})
		require.EqualError(t, err, "some-get-error")
	})

}

func TestBookController_List(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	bookSvc := mock.NewMockBookService(ctrl)
	bookCntrl := controller.BookController{
		BookService: bookSvc,
	}

	t.Run("When repo success", func(t *testing.T) {
		bookSvc.EXPECT().List(gomock.Any()).Return([]*repository.Book{
			&repository.Book{ID: 1, Title: "title1", Author: "author1"},
			&repository.Book{ID: 2, Title: "title2", Author: "author2"},
		}, nil)
		rr, err := echokit.DoGET(bookCntrl.List, "/", nil)
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, rr.Code)
		require.Equal(t, "[{\"id\":1,\"title\":\"title1\",\"author\":\"author1\"},{\"id\":2,\"title\":\"title2\",\"author\":\"author2\"}]\n", rr.Body.String())
	})

	t.Run("When repo error", func(t *testing.T) {
		bookSvc.EXPECT().List(gomock.Any()).Return(nil, fmt.Errorf("some-list-error"))
		_, err := echokit.DoGET(bookCntrl.List, "/", nil)
		require.EqualError(t, err, "some-list-error")
	})
}

func TestBookController_Create(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	bookSvc := mock.NewMockBookService(ctrl)
	bookController := controller.BookController{
		BookService: bookSvc,
	}

	t.Run("When invalid message request", func(t *testing.T) {
		rr, err := echokit.DoPOST(bookController.Create, "/", `{}`)
		require.NoError(t, err)
		require.Equal(t, http.StatusBadRequest, rr.Code)
		require.Equal(t, "{\"message\":\"Invalid Message\"}\n", rr.Body.String())
	})

	t.Run("When invalid json format", func(t *testing.T) {
		_, err := echokit.DoPOST(bookController.Create, "/", `invalid}`)
		require.EqualError(t, err, `code=400, message=Syntax error: offset=1, error=invalid character 'i' looking for beginning of value`)
	})

	t.Run("When error", func(t *testing.T) {
		bookSvc.EXPECT().Insert(gomock.Any(), gomock.Any()).Return(int64(0), fmt.Errorf("some-insert-error"))
		_, err := echokit.DoPOST(bookController.Create, "/", `{"author":"some-author", "title":"some-title"}`)
		require.EqualError(t, err, "some-insert-error")
	})

	t.Run("When success", func(t *testing.T) {
		bookSvc.EXPECT().Insert(gomock.Any(), gomock.Any()).Return(int64(99), nil)
		rr, err := echokit.DoPOST(bookController.Create, "/", `{"author":"some-author", "title":"some-title"}`)
		require.NoError(t, err)
		require.Equal(t, http.StatusCreated, rr.Code)
		require.Equal(t, "{\"message\":\"Success insert new record #99\"}\n", rr.Body.String())
	})
}

func TestBookController_Delete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	bookSvc := mock.NewMockBookService(ctrl)
	bookCntrl := controller.BookController{
		BookService: bookSvc,
	}
	t.Run("When invalid ID", func(t *testing.T) {
		rr, err := echokit.DoDELETE(bookCntrl.Delete, "/", map[string]string{
			"id": "invalid",
		})
		require.NoError(t, err)
		require.Equal(t, http.StatusBadRequest, rr.Code)
	})
	t.Run("When return success", func(t *testing.T) {
		bookSvc.EXPECT().Delete(gomock.Any(), int64(1)).Return(nil)
		rr, err := echokit.DoDELETE(bookCntrl.Delete, "/", map[string]string{
			"id": "1",
		})
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, rr.Code)
		require.Equal(t, "{\"message\":\"Delete #1 done\"}\n", rr.Body.String())
	})
	t.Run("When error", func(t *testing.T) {
		bookSvc.EXPECT().Delete(gomock.Any(), int64(2)).Return(fmt.Errorf("some-delete-error"))
		_, err := echokit.DoDELETE(bookCntrl.Delete, "/", map[string]string{
			"id": "2",
		})
		require.EqualError(t, err, "some-delete-error")
	})
}

func TestBookController_Update(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	bookSvc := mock.NewMockBookService(ctrl)
	bookCntrl := controller.BookController{
		BookService: bookSvc,
	}
	t.Run("When invalid message request", func(t *testing.T) {
		rr, err := echokit.DoPUT(bookCntrl.Update, "/", `{}`)
		require.NoError(t, err)
		require.Equal(t, http.StatusBadRequest, rr.Code)
		require.Equal(t, "{\"message\":\"Invalid ID\"}\n", rr.Body.String())
	})

	t.Run("When invalid message request", func(t *testing.T) {
		rr, err := echokit.DoPUT(bookCntrl.Update, "/", `{"id": 1}`)
		require.NoError(t, err)
		require.Equal(t, http.StatusBadRequest, rr.Code)
		require.Equal(t, "{\"message\":\"Invalid Message\"}\n", rr.Body.String())
	})

	t.Run("When invalid json format", func(t *testing.T) {
		_, err := echokit.DoPUT(bookCntrl.Update, "/", `invalid}`)
		require.EqualError(t, err, `code=400, message=Syntax error: offset=1, error=invalid character 'i' looking for beginning of value`)
	})

	t.Run("When error", func(t *testing.T) {
		bookSvc.EXPECT().Update(gomock.Any(), gomock.Any()).Return(fmt.Errorf("some-update-error"))
		_, err := echokit.DoPUT(bookCntrl.Update, "/", `{"id": 1,"author":"some-author", "title":"some-title"}`)
		require.EqualError(t, err, "some-update-error")
	})

	t.Run("When success", func(t *testing.T) {
		bookSvc.EXPECT().Update(gomock.Any(), gomock.Any()).Return(nil)
		rr, err := echokit.DoPUT(bookCntrl.Update, "/", `{"id": 1, "author":"some-author", "title":"some-title"}`)
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, rr.Code)
		require.Equal(t, "{\"message\":\"Update #1 success\"}\n", rr.Body.String())
	})
}
