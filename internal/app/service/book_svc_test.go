package service_test

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"

	"github.com/typical-go/typical-rest-server/internal/app/repo"
	"github.com/typical-go/typical-rest-server/internal/app/service"
	"github.com/typical-go/typical-rest-server/internal/generated/mock/app/repo_mock"
	"github.com/typical-go/typical-rest-server/pkg/sqkit"
)

type OnMockBookRepo func(mockRepo *repo_mock.MockBookRepo)

func createBookSvc(t *testing.T, onMockBookRepo OnMockBookRepo) (*service.BookSvcImpl, *gomock.Controller) {
	mock := gomock.NewController(t)
	mockRepo := repo_mock.NewMockBookRepo(mock)
	if onMockBookRepo != nil {
		onMockBookRepo(mockRepo)
	}

	return &service.BookSvcImpl{
		Repo: mockRepo,
	}, mock
}

func TestBookSvc_Create(t *testing.T) {
	testcases := []struct {
		testName       string
		OnMockBookRepo OnMockBookRepo
		book           *repo.Book
		expected       *repo.Book
		expectedErr    string
	}{
		{
			testName:    "title is missing",
			book:        &repo.Book{},
			expectedErr: "code=422, message=title must be filled",
		},
		{
			testName:    "author is missing",
			book:        &repo.Book{Title: "some-title"},
			expectedErr: "code=422, message=author must be filled",
		},
		{
			testName:    "create error",
			book:        &repo.Book{Author: "some-author", Title: "some-title"},
			expectedErr: "create-error",
			OnMockBookRepo: func(mockRepo *repo_mock.MockBookRepo) {
				mockRepo.EXPECT().
					Insert(gomock.Any(), &repo.Book{Author: "some-author", Title: "some-title"}).
					Return(int64(-1), errors.New("create-error"))
			},
		},
		{
			testName:    "find error",
			book:        &repo.Book{Author: "some-author", Title: "some-title"},
			expectedErr: "find-error",
			OnMockBookRepo: func(mockRepo *repo_mock.MockBookRepo) {
				mockRepo.EXPECT().
					Insert(gomock.Any(), &repo.Book{Author: "some-author", Title: "some-title"}).
					Return(int64(1), nil)
				mockRepo.EXPECT().
					Find(gomock.Any(), sqkit.Eq{"id": int64(1)}).
					Return(nil, errors.New("find-error"))
			},
		},
		{
			book: &repo.Book{
				Author: "some-author",
				Title:  "some-title",
			},
			expected: &repo.Book{Author: "some-author", Title: "some-title"},
			OnMockBookRepo: func(mockRepo *repo_mock.MockBookRepo) {
				mockRepo.EXPECT().
					Insert(gomock.Any(), &repo.Book{Author: "some-author", Title: "some-title"}).
					Return(int64(1), nil)
				mockRepo.EXPECT().
					Find(gomock.Any(), sqkit.Eq{"id": int64(1)}).
					Return([]*repo.Book{{Author: "some-author", Title: "some-title"}}, nil)
			},
		},
	}

	for _, tt := range testcases {
		t.Run(tt.testName, func(t *testing.T) {
			svc, mock := createBookSvc(t, tt.OnMockBookRepo)
			defer mock.Finish()

			id, err := svc.Create(context.Background(), tt.book)
			if tt.expectedErr != "" {
				require.EqualError(t, err, tt.expectedErr)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.expected, id)
			}
		})
	}
}

func TestBookSvc_FindOne(t *testing.T) {
	testcases := []struct {
		testName       string
		OnMockBookRepo OnMockBookRepo
		paramID        string
		expected       *repo.Book
		expectedErr    string
	}{
		{
			paramID: "1",
			OnMockBookRepo: func(mockRepo *repo_mock.MockBookRepo) {
				mockRepo.EXPECT().
					Find(gomock.Any(), sqkit.Eq{"id": int64(1)}).
					Return(nil, errors.New("some-error"))
			},
			expectedErr: "some-error",
		},
		{
			paramID: "1",
			OnMockBookRepo: func(mockRepo *repo_mock.MockBookRepo) {
				mockRepo.EXPECT().
					Find(gomock.Any(), sqkit.Eq{"id": int64(1)}).
					Return([]*repo.Book{{ID: 1, Title: "some-title"}}, nil)
			},
			expected: &repo.Book{ID: 1, Title: "some-title"},
		},
		{
			paramID: "1",
			OnMockBookRepo: func(mockRepo *repo_mock.MockBookRepo) {
				mockRepo.EXPECT().
					Find(gomock.Any(), sqkit.Eq{"id": int64(1)}).
					Return([]*repo.Book{}, nil)
			},
			expectedErr: "code=404, message=Not Found",
		},
	}
	for _, tt := range testcases {
		t.Run(tt.testName, func(t *testing.T) {
			svc, mock := createBookSvc(t, tt.OnMockBookRepo)
			defer mock.Finish()

			book, err := svc.FindOne(context.Background(), tt.paramID)
			if tt.expectedErr != "" {
				require.EqualError(t, err, tt.expectedErr)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.expected, book)
			}
		})
	}
}

func TestBookSvc_Find(t *testing.T) {
	testcases := []struct {
		testName       string
		OnMockBookRepo OnMockBookRepo
		req            *service.FindBookReq
		expected       *service.FindBookResp
		expectedErr    string
	}{
		{
			OnMockBookRepo: func(mockRepo *repo_mock.MockBookRepo) {
				mockRepo.EXPECT().Count(gomock.Any()).Return(int64(10), nil)
				mockRepo.EXPECT().
					Find(gomock.Any(), &sqkit.OffsetPagination{}).
					Return([]*repo.Book{
						{ID: 1, Title: "title1", Author: "author1"},
						{ID: 2, Title: "title2", Author: "author2"},
					}, nil)
			},
			req: &service.FindBookReq{},
			expected: &service.FindBookResp{
				Books: []*repo.Book{
					{ID: 1, Title: "title1", Author: "author1"},
					{ID: 2, Title: "title2", Author: "author2"},
				},
				TotalCount: "10",
			},
		},
		{
			testName: "count error",
			OnMockBookRepo: func(mockRepo *repo_mock.MockBookRepo) {
				mockRepo.EXPECT().
					Count(gomock.Any()).
					Return(int64(-1), errors.New("count-error"))
			},
			req:         &service.FindBookReq{Limit: 20, Offset: 10, Sort: "title,created_at"},
			expectedErr: "count-error",
		},
		{
			testName: "find error",
			OnMockBookRepo: func(mockRepo *repo_mock.MockBookRepo) {
				mockRepo.EXPECT().
					Count(gomock.Any()).
					Return(int64(10), nil)
				mockRepo.EXPECT().
					Find(gomock.Any(), &sqkit.OffsetPagination{Limit: 20, Offset: 10}, sqkit.Sorts{"title", "created_at"}).
					Return(nil, errors.New("find-error"))
			},
			req:         &service.FindBookReq{Limit: 20, Offset: 10, Sort: "title,created_at"},
			expectedErr: "find-error",
		},
	}
	for _, tt := range testcases {
		t.Run(tt.testName, func(t *testing.T) {
			svc, mock := createBookSvc(t, tt.OnMockBookRepo)
			defer mock.Finish()

			books, err := svc.Find(context.Background(), tt.req)
			if tt.expectedErr != "" {
				require.EqualError(t, err, tt.expectedErr)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.expected, books)
			}
		})
	}
}

func TestBookSvc_Delete(t *testing.T) {
	testcases := []struct {
		testName       string
		OnMockBookRepo OnMockBookRepo
		paramID        string
		expectedErr    string
	}{
		{
			paramID:     "1",
			expectedErr: `some-error`,
			OnMockBookRepo: func(mockRepo *repo_mock.MockBookRepo) {
				mockRepo.EXPECT().
					Delete(gomock.Any(), sqkit.Eq{repo.BookTable.ID: int64(1)}).
					Return(int64(0), errors.New("some-error"))
			},
		},
		{
			paramID: "1",
			OnMockBookRepo: func(mockRepo *repo_mock.MockBookRepo) {
				mockRepo.EXPECT().
					Delete(gomock.Any(), sqkit.Eq{repo.BookTable.ID: int64(1)}).
					Return(int64(1), nil)
			},
		},
		{
			testName: "success even if no affected row (idempotent)",
			paramID:  "1",
			OnMockBookRepo: func(mockRepo *repo_mock.MockBookRepo) {
				mockRepo.EXPECT().
					Delete(gomock.Any(), sqkit.Eq{repo.BookTable.ID: int64(1)}).
					Return(int64(0), nil)
			},
		},
	}
	for _, tt := range testcases {
		t.Run(tt.testName, func(t *testing.T) {
			svc, mock := createBookSvc(t, tt.OnMockBookRepo)
			defer mock.Finish()

			err := svc.Delete(context.Background(), tt.paramID)
			if tt.expectedErr != "" {
				require.EqualError(t, err, tt.expectedErr)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestBookSvc_Update(t *testing.T) {
	testcases := []struct {
		testName       string
		OnMockBookRepo OnMockBookRepo
		paramID        string
		book           *repo.Book
		expected       *repo.Book
		expectedErr    string
	}{
		{
			testName:    "missing title",
			paramID:     "1",
			book:        &repo.Book{},
			expectedErr: "code=422, message=title must be filled",
		},
		{
			testName:    "missing author",
			paramID:     "1",
			book:        &repo.Book{Title: "some-title"},
			expectedErr: "code=422, message=author must be filled",
		},
		{
			testName:    "update error",
			paramID:     "1",
			book:        &repo.Book{Author: "some-author", Title: "some-title"},
			expectedErr: "update error",
			OnMockBookRepo: func(mockRepo *repo_mock.MockBookRepo) {
				mockRepo.EXPECT().
					Find(gomock.Any(), sqkit.Eq{"id": int64(1)}).
					Return([]*repo.Book{{ID: 1, Title: "some-title"}}, nil)
				mockRepo.EXPECT().
					Update(gomock.Any(), &repo.Book{Author: "some-author", Title: "some-title"}, sqkit.Eq{"id": int64(1)}).
					Return(int64(-1), errors.New("update error"))
			},
		},
		{
			testName:    "nothing to update",
			paramID:     "1",
			book:        &repo.Book{Author: "some-author", Title: "some-title"},
			expectedErr: "no affected row",
			OnMockBookRepo: func(mockRepo *repo_mock.MockBookRepo) {
				mockRepo.EXPECT().
					Find(gomock.Any(), sqkit.Eq{"id": int64(1)}).
					Return([]*repo.Book{{ID: 1, Title: "some-title"}}, nil)
				mockRepo.EXPECT().
					Update(gomock.Any(), &repo.Book{Author: "some-author", Title: "some-title"}, sqkit.Eq{"id": int64(1)}).
					Return(int64(0), nil)
			},
		},
		{
			testName:    "find error before update",
			paramID:     "1",
			book:        &repo.Book{Author: "some-author", Title: "some-title"},
			expectedErr: "find-error",
			OnMockBookRepo: func(mockRepo *repo_mock.MockBookRepo) {
				mockRepo.EXPECT().
					Find(gomock.Any(), sqkit.Eq{"id": int64(1)}).
					Return(nil, errors.New("find-error"))
			},
		},
		{
			testName:    "find error after update",
			paramID:     "1",
			book:        &repo.Book{Author: "some-author", Title: "some-title"},
			expectedErr: "find-error",
			OnMockBookRepo: func(mockRepo *repo_mock.MockBookRepo) {
				mockRepo.EXPECT().
					Find(gomock.Any(), sqkit.Eq{"id": int64(1)}).
					Return([]*repo.Book{{ID: 1, Title: "some-title"}}, nil)
				mockRepo.EXPECT().
					Update(gomock.Any(), &repo.Book{Author: "some-author", Title: "some-title"}, sqkit.Eq{"id": int64(1)}).
					Return(int64(1), nil)
				mockRepo.EXPECT().
					Find(gomock.Any(), sqkit.Eq{"id": int64(1)}).
					Return(nil, errors.New("find-error"))
			},
		},
	}
	for _, tt := range testcases {
		t.Run(tt.testName, func(t *testing.T) {
			svc, mock := createBookSvc(t, tt.OnMockBookRepo)
			defer mock.Finish()

			book, err := svc.Update(context.Background(), tt.paramID, tt.book)
			if tt.expectedErr != "" {
				require.EqualError(t, err, tt.expectedErr)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.expected, book)
			}
		})
	}
}

func TestBookSvc_Patch(t *testing.T) {
	testcases := []struct {
		testName       string
		OnMockBookRepo OnMockBookRepo
		paramID        string
		book           *repo.Book
		expected       *repo.Book
		expectedErr    string
	}{
		{
			testName:    "patch error",
			paramID:     "1",
			book:        &repo.Book{Author: "some-author", Title: "some-title"},
			expectedErr: "patch-error",
			OnMockBookRepo: func(mockRepo *repo_mock.MockBookRepo) {
				mockRepo.EXPECT().
					Find(gomock.Any(), sqkit.Eq{"id": int64(1)}).
					Return([]*repo.Book{{ID: 1, Title: "some-title"}}, nil)
				mockRepo.EXPECT().
					Patch(gomock.Any(), &repo.Book{Author: "some-author", Title: "some-title"}, sqkit.Eq{"id": int64(1)}).
					Return(int64(-1), errors.New("patch-error"))
			},
		},
		{
			testName:    "patch error",
			paramID:     "1",
			book:        &repo.Book{Author: "some-author", Title: "some-title"},
			expectedErr: "no affected row",
			OnMockBookRepo: func(mockRepo *repo_mock.MockBookRepo) {
				mockRepo.EXPECT().
					Find(gomock.Any(), sqkit.Eq{"id": int64(1)}).
					Return([]*repo.Book{{ID: 1, Title: "some-title"}}, nil)
				mockRepo.EXPECT().
					Patch(gomock.Any(), &repo.Book{Author: "some-author", Title: "some-title"}, sqkit.Eq{"id": int64(1)}).
					Return(int64(0), nil)
			},
		},
		{
			testName:    "find error before update",
			paramID:     "1",
			book:        &repo.Book{Author: "some-author", Title: "some-title"},
			expectedErr: "find-error",
			OnMockBookRepo: func(mockRepo *repo_mock.MockBookRepo) {
				mockRepo.EXPECT().
					Find(gomock.Any(), sqkit.Eq{"id": int64(1)}).
					Return(nil, errors.New("find-error"))
			},
		},
		{
			testName:    "find error after update",
			paramID:     "1",
			book:        &repo.Book{Author: "some-author", Title: "some-title"},
			expectedErr: "find-error",
			OnMockBookRepo: func(mockRepo *repo_mock.MockBookRepo) {
				mockRepo.EXPECT().
					Find(gomock.Any(), sqkit.Eq{"id": int64(1)}).
					Return([]*repo.Book{{ID: 1, Title: "some-title"}}, nil)
				mockRepo.EXPECT().
					Patch(gomock.Any(), &repo.Book{Author: "some-author", Title: "some-title"}, sqkit.Eq{"id": int64(1)}).
					Return(int64(1), nil)
				mockRepo.EXPECT().
					Find(gomock.Any(), sqkit.Eq{"id": int64(1)}).
					Return(nil, errors.New("find-error"))
			},
		},
		{
			paramID:  "1",
			book:     &repo.Book{Author: "some-author", Title: "some-title"},
			expected: &repo.Book{Author: "some-author", Title: "some-title"},
			OnMockBookRepo: func(mockRepo *repo_mock.MockBookRepo) {
				mockRepo.EXPECT().
					Find(gomock.Any(), sqkit.Eq{"id": int64(1)}).
					Return([]*repo.Book{{ID: 1, Title: "some-title"}}, nil)
				mockRepo.EXPECT().
					Patch(gomock.Any(), &repo.Book{Author: "some-author", Title: "some-title"}, sqkit.Eq{"id": int64(1)}).
					Return(int64(1), nil)
				mockRepo.EXPECT().
					Find(gomock.Any(), sqkit.Eq{"id": int64(1)}).
					Return([]*repo.Book{{Author: "some-author", Title: "some-title"}}, nil)
			},
		},
	}
	for _, tt := range testcases {
		t.Run(tt.testName, func(t *testing.T) {
			svc, mock := createBookSvc(t, tt.OnMockBookRepo)
			defer mock.Finish()
			book, err := svc.Patch(context.Background(), tt.paramID, tt.book)
			if tt.expectedErr != "" {
				require.EqualError(t, err, tt.expectedErr)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.expected, book)
			}
		})
	}
}
