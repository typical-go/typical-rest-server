package cntrl_test

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/typical-go/typical-rest-server/app/cntrl"
	"github.com/typical-go/typical-rest-server/app/helper/echokit"

	"github.com/typical-go/typical-rest-server/app/repo"
	"github.com/typical-go/typical-rest-server/mock"

	"github.com/stretchr/testify/require"
)

func TestBookController_NoRepository(t *testing.T) {
	bookController := cntrl.NewBookController(nil)
	ctx, _ := echokit.RequestGET("/")

	require.EqualError(t, bookController.Get(ctx), "BookRepository is missing")
	require.EqualError(t, bookController.Delete(ctx), "BookRepository is missing")
	require.EqualError(t, bookController.Create(ctx), "BookRepository is missing")
	require.EqualError(t, bookController.Update(ctx), "BookRepository is missing")
	require.EqualError(t, bookController.List(ctx), "BookRepository is missing")
}

func TestBookController(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	bookR := mock.NewMockBookRepository(ctrl)
	bookController := cntrl.NewBookController(bookR)

	t.Run("Get", func(t *testing.T) {
		t.Run("When invalid ID", func(t *testing.T) {
			ctx, rr := echokit.RequestGETWithParam("/", map[string]string{
				"id": "invalid",
			})
			err := bookController.Get(ctx)
			require.NoError(t, err)
			require.Equal(t, http.StatusBadRequest, rr.Code)
		})

		t.Run("When return success", func(t *testing.T) {
			bookR.EXPECT().Find(int64(1)).Return(&repo.Book{ID: 1, Title: "title1", Author: "author1"}, nil)

			ctx, rr := echokit.RequestGETWithParam("/", map[string]string{
				"id": "1",
			})

			err := bookController.Get(ctx)
			require.NoError(t, err)
			require.Equal(t, http.StatusOK, rr.Code)
			require.Equal(t, "{\"id\":1,\"title\":\"title1\",\"author\":\"author1\"}\n", rr.Body.String())
		})

		t.Run("When entity not found", func(t *testing.T) {
			bookR.EXPECT().Find(int64(3)).Return(nil, nil)

			ctx, rr := echokit.RequestGETWithParam("/", map[string]string{
				"id": "3",
			})

			err := bookController.Get(ctx)
			require.NoError(t, err)
			require.Equal(t, http.StatusNotFound, rr.Code)
			require.Equal(t, "{\"message\":\"book #3 not found\"}\n", rr.Body.String())
		})

		t.Run("When return error", func(t *testing.T) {
			bookR.EXPECT().Find(int64(2)).Return(nil, fmt.Errorf("some-get-error"))

			ctx, _ := echokit.RequestGETWithParam("/", map[string]string{
				"id": "2",
			})

			err := bookController.Get(ctx)
			require.EqualError(t, err, "some-get-error")
		})
	})

	t.Run("List", func(t *testing.T) {
		t.Run("When repo success", func(t *testing.T) {
			bookR.EXPECT().List().Return(
				[]*repo.Book{
					&repo.Book{ID: 1, Title: "title1", Author: "author1"},
					&repo.Book{ID: 2, Title: "title2", Author: "author2"},
				}, nil)

			ctx, rr := echokit.RequestGET("/")
			err := bookController.List(ctx)
			require.NoError(t, err)
			require.Equal(t, http.StatusOK, rr.Code)
			require.Equal(t, "[{\"id\":1,\"title\":\"title1\",\"author\":\"author1\"},{\"id\":2,\"title\":\"title2\",\"author\":\"author2\"}]\n", rr.Body.String())
		})

		t.Run("When repo error", func(t *testing.T) {
			bookR.EXPECT().List().Return(nil, fmt.Errorf("some-list-error"))

			ctx, _ := echokit.RequestGET("/")
			err := bookController.List(ctx)
			require.EqualError(t, err, "some-list-error")
		})
	})

	t.Run("Create", func(t *testing.T) {
		t.Run("When invalid message request", func(t *testing.T) {
			ctx, rr := echokit.RequestPOST("/", `{}`)
			err := bookController.Create(ctx)
			require.NoError(t, err)
			require.Equal(t, http.StatusBadRequest, rr.Code)
			require.Equal(t, "{\"message\":\"Invalid Message\"}\n", rr.Body.String())
		})

		t.Run("When invalid json format", func(t *testing.T) {
			ctx, _ := echokit.RequestPOST("/", `invalid-json}`)
			err := bookController.Create(ctx)
			require.EqualError(t, err, `code=400, message=Syntax error: offset=1, error=invalid character 'i' looking for beginning of value`)
		})

		t.Run("When error", func(t *testing.T) {
			bookR.EXPECT().Insert(gomock.Any()).Return(int64(0), fmt.Errorf("some-insert-error"))

			ctx, _ := echokit.RequestPOST("/", `{"author":"some-author", "title":"some-title"}`)
			err := bookController.Create(ctx)
			require.EqualError(t, err, "some-insert-error")
		})

		t.Run("When success", func(t *testing.T) {
			bookR.EXPECT().Insert(gomock.Any()).Return(int64(99), nil)

			ctx, rr := echokit.RequestPOST("/", `{"author":"some-author", "title":"some-title"}`)
			err := bookController.Create(ctx)
			require.NoError(t, err)
			require.Equal(t, http.StatusCreated, rr.Code)
			require.Equal(t, "{\"message\":\"Success insert new record #99\"}\n", rr.Body.String())
		})
	})

	t.Run("Delete", func(t *testing.T) {
		t.Run("When invalid ID", func(t *testing.T) {
			ctx, rr := echokit.RequestGETWithParam("/", map[string]string{
				"id": "invalid",
			})
			err := bookController.Delete(ctx)
			require.NoError(t, err)
			require.Equal(t, http.StatusBadRequest, rr.Code)
		})

		t.Run("When return success", func(t *testing.T) {
			bookR.EXPECT().Delete(int64(1)).Return(nil)

			ctx, rr := echokit.RequestGETWithParam("/", map[string]string{
				"id": "1",
			})

			err := bookController.Delete(ctx)
			require.NoError(t, err)
			require.Equal(t, http.StatusOK, rr.Code)
			require.Equal(t, "{\"message\":\"Delete #1 done\"}\n", rr.Body.String())
		})

		t.Run("When error", func(t *testing.T) {
			bookR.EXPECT().Delete(int64(2)).Return(fmt.Errorf("some-delete-error"))

			ctx, _ := echokit.RequestGETWithParam("/", map[string]string{
				"id": "2",
			})

			err := bookController.Delete(ctx)
			require.EqualError(t, err, "some-delete-error")
		})
	})

	t.Run("Updte", func(t *testing.T) {
		t.Run("When invalid message request", func(t *testing.T) {
			ctx, rr := echokit.RequestPOST("/", `{}`)
			err := bookController.Update(ctx)
			require.NoError(t, err)
			require.Equal(t, http.StatusBadRequest, rr.Code)
			require.Equal(t, "{\"message\":\"Invalid ID\"}\n", rr.Body.String())
		})

		t.Run("When invalid message request", func(t *testing.T) {
			ctx, rr := echokit.RequestPOST("/", `{"id": 1}`)
			err := bookController.Update(ctx)
			require.NoError(t, err)
			require.Equal(t, http.StatusBadRequest, rr.Code)
			require.Equal(t, "{\"message\":\"Invalid Message\"}\n", rr.Body.String())
		})

		t.Run("When invalid json format", func(t *testing.T) {
			ctx, _ := echokit.RequestPOST("/", `invalid-json}`)
			err := bookController.Update(ctx)
			require.EqualError(t, err, `code=400, message=Syntax error: offset=1, error=invalid character 'i' looking for beginning of value`)
		})

		t.Run("When error", func(t *testing.T) {
			bookR.EXPECT().Update(gomock.Any()).Return(fmt.Errorf("some-update-error"))

			ctx, _ := echokit.RequestPOST("/", `{"id": 1,"author":"some-author", "title":"some-title"}`)
			err := bookController.Update(ctx)
			require.EqualError(t, err, "some-update-error")
		})

		t.Run("When success", func(t *testing.T) {
			bookR.EXPECT().Update(gomock.Any()).Return(nil)

			ctx, rr := echokit.RequestPOST("/", `{"id": 1, "author":"some-author", "title":"some-title"}`)
			err := bookController.Update(ctx)
			require.NoError(t, err)
			require.Equal(t, http.StatusOK, rr.Code)
			require.Equal(t, "{\"message\":\"Update success\"}\n", rr.Body.String())
		})
	})
}
