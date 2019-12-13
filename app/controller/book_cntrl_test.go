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
	bookCntrl := controller.BookCntrl{
		BookService: bookSvc,
	}
	t.Run("GIVEN invalid id", func(t *testing.T) {
		_, err := echokit.DoGET(bookCntrl.Get, "/", map[string]string{
			"id": "invalid",
		})
		require.EqualError(t, err, "code=400, message=Invalid ID")
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
	t.Run("WHEN repository not found", func(t *testing.T) {
		bookSvc.EXPECT().Find(gomock.Any(), int64(3)).Return(nil, nil)
		_, err := echokit.DoGET(bookCntrl.Get, "/", map[string]string{
			"id": "3",
		})
		require.EqualError(t, err, "code=404, message=Book#3 not found")
	})
	t.Run("WHEN return error", func(t *testing.T) {
		bookSvc.EXPECT().Find(gomock.Any(), int64(2)).Return(nil, fmt.Errorf("some-get-error"))
		_, err := echokit.DoGET(bookCntrl.Get, "/", map[string]string{
			"id": "2",
		})
		require.EqualError(t, err, "code=500, message=some-get-error")
	})

}

func TestBookController_List(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	bookSvc := mock.NewMockBookService(ctrl)
	bookCntrl := controller.BookCntrl{
		BookService: bookSvc,
	}
	t.Run("WHEN repo success", func(t *testing.T) {
		bookSvc.EXPECT().List(gomock.Any()).Return([]*repository.Book{
			&repository.Book{ID: 1, Title: "title1", Author: "author1"},
			&repository.Book{ID: 2, Title: "title2", Author: "author2"},
		}, nil)
		rr, err := echokit.DoGET(bookCntrl.List, "/", nil)
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, rr.Code)
		require.Equal(t, "[{\"id\":1,\"title\":\"title1\",\"author\":\"author1\"},{\"id\":2,\"title\":\"title2\",\"author\":\"author2\"}]\n", rr.Body.String())
	})
	t.Run("WHEN repo error", func(t *testing.T) {
		bookSvc.EXPECT().List(gomock.Any()).Return(nil, fmt.Errorf("some-list-error"))
		_, err := echokit.DoGET(bookCntrl.List, "/", nil)
		require.EqualError(t, err, "code=500, message=some-list-error")
	})
}

func TestBookController_Create(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	bookSvc := mock.NewMockBookService(ctrl)
	bookController := controller.BookCntrl{
		BookService: bookSvc,
	}
	t.Run("WHEN invalid message request", func(t *testing.T) {
		_, err := echokit.DoPOST(bookController.Create, "/", `{}`)
		require.EqualError(t, err, "code=400, message=Key: 'Book.Title' Error:Field validation for 'Title' failed on the 'required' tag\nKey: 'Book.Author' Error:Field validation for 'Author' failed on the 'required' tag")
	})
	t.Run("WHEN invalid json format", func(t *testing.T) {
		_, err := echokit.DoPOST(bookController.Create, "/", `invalid}`)
		require.EqualError(t, err, `code=400, message=Syntax error: offset=1, error=invalid character 'i' looking for beginning of value`)
	})
	t.Run("WHEN error", func(t *testing.T) {
		bookSvc.EXPECT().Insert(gomock.Any(), gomock.Any()).Return(int64(0), fmt.Errorf("some-insert-error"))
		_, err := echokit.DoPOST(bookController.Create, "/", `{"author":"some-author", "title":"some-title"}`)
		require.EqualError(t, err, "code=422, message=some-insert-error")
	})
	t.Run("WHEN success", func(t *testing.T) {
		bookSvc.EXPECT().Insert(gomock.Any(), gomock.Any()).Return(int64(99), nil)
		rr, err := echokit.DoPOST(bookController.Create, "/", `{"author":"some-author", "title":"some-title"}`)
		require.NoError(t, err)
		require.Equal(t, http.StatusCreated, rr.Code)
		require.Equal(t, "{\"message\":\"Success insert new book #99\"}\n", rr.Body.String())
	})
}

func TestBookController_Delete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	bookSvc := mock.NewMockBookService(ctrl)
	bookCntrl := controller.BookCntrl{
		BookService: bookSvc,
	}
	t.Run("WHEN invalid ID", func(t *testing.T) {
		_, err := echokit.DoDELETE(bookCntrl.Delete, "/", map[string]string{
			"id": "invalid",
		})
		require.EqualError(t, err, "code=400, message=Invalid ID")
	})
	t.Run("WHEN return success", func(t *testing.T) {
		bookSvc.EXPECT().Delete(gomock.Any(), int64(1)).Return(nil)
		rr, err := echokit.DoDELETE(bookCntrl.Delete, "/", map[string]string{
			"id": "1",
		})
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, rr.Code)
		require.Equal(t, "{\"message\":\"Success delete book #1\"}\n", rr.Body.String())
	})
	t.Run("WHEN error", func(t *testing.T) {
		bookSvc.EXPECT().Delete(gomock.Any(), int64(2)).Return(fmt.Errorf("some-delete-error"))
		_, err := echokit.DoDELETE(bookCntrl.Delete, "/", map[string]string{
			"id": "2",
		})
		require.EqualError(t, err, "code=500, message=some-delete-error")
	})
}

func TestBookController_Update(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	bookSvc := mock.NewMockBookService(ctrl)
	bookCntrl := controller.BookCntrl{
		BookService: bookSvc,
	}
	t.Run("When invalid message request", func(t *testing.T) {
		_, err := echokit.DoPUT(bookCntrl.Update, "/", `{}`)
		require.EqualError(t, err, "code=400, message=Invalid ID")
	})
	t.Run("When invalid message request", func(t *testing.T) {
		_, err := echokit.DoPUT(bookCntrl.Update, "/", `{"id": 1}`)
		require.EqualError(t, err, "code=400, message=Key: 'Book.Title' Error:Field validation for 'Title' failed on the 'required' tag\nKey: 'Book.Author' Error:Field validation for 'Author' failed on the 'required' tag")
	})
	t.Run("When invalid json format", func(t *testing.T) {
		_, err := echokit.DoPUT(bookCntrl.Update, "/", `invalid}`)
		require.EqualError(t, err, `code=400, message=Syntax error: offset=1, error=invalid character 'i' looking for beginning of value`)
	})
	t.Run("When error", func(t *testing.T) {
		bookSvc.EXPECT().Update(gomock.Any(), gomock.Any()).Return(fmt.Errorf("some-update-error"))
		_, err := echokit.DoPUT(bookCntrl.Update, "/", `{"id": 1,"author":"some-author", "title":"some-title"}`)
		require.EqualError(t, err, "code=500, message=some-update-error")
	})
	t.Run("When success", func(t *testing.T) {
		bookSvc.EXPECT().Update(gomock.Any(), gomock.Any()).Return(nil)
		rr, err := echokit.DoPUT(bookCntrl.Update, "/", `{"id": 1, "author":"some-author", "title":"some-title"}`)
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, rr.Code)
		require.Equal(t, "{\"message\":\"Success update book #1\"}\n", rr.Body.String())
	})
}
