package service_test

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"github.com/typical-go/typical-rest-server/internal/app/data_access/librarydb"
	"github.com/typical-go/typical-rest-server/internal/app/data_access/librarydb_mock"
	"github.com/typical-go/typical-rest-server/internal/app/domain/library/service"
	"github.com/typical-go/typical-rest-server/pkg/dbkit"
)

type bookSvcFn func(mockRepo *librarydb_mock.MockBookRepo)

func createBookSvc(t *testing.T, fn bookSvcFn) (*service.BookSvcImpl, *gomock.Controller) {
	mock := gomock.NewController(t)
	mockRepo := librarydb_mock.NewMockBookRepo(mock)
	if fn != nil {
		fn(mockRepo)
	}

	return &service.BookSvcImpl{
		BookRepo: mockRepo,
	}, mock
}

func TestBookSvc_Create(t *testing.T) {
	testcases := []struct {
		testName    string
		bookSvcFn   bookSvcFn
		book        *librarydb.Book
		expected    *librarydb.Book
		expectedErr string
	}{
		{
			testName:    "validation error",
			book:        &librarydb.Book{},
			expectedErr: "Key: 'Book.Title' Error:Field validation for 'Title' failed on the 'required' tag\nKey: 'Book.Author' Error:Field validation for 'Author' failed on the 'required' tag",
		},
		{
			testName:    "create error",
			book:        &librarydb.Book{Author: "some-author", Title: "some-title"},
			expectedErr: "create-error",
			bookSvcFn: func(mockRepo *librarydb_mock.MockBookRepo) {
				mockRepo.EXPECT().
					Create(gomock.Any(), &librarydb.Book{Author: "some-author", Title: "some-title"}).
					Return(int64(-1), errors.New("create-error"))
			},
		},
		{
			testName:    "Find error",
			book:        &librarydb.Book{Author: "some-author", Title: "some-title"},
			expectedErr: "Find-error",
			bookSvcFn: func(mockRepo *librarydb_mock.MockBookRepo) {
				mockRepo.EXPECT().
					Create(gomock.Any(), &librarydb.Book{Author: "some-author", Title: "some-title"}).
					Return(int64(1), nil)
				mockRepo.EXPECT().
					Find(gomock.Any(), dbkit.Equal("id", int64(1))).
					Return(nil, errors.New("Find-error"))
			},
		},
		{
			book: &librarydb.Book{
				Author: "some-author",
				Title:  "some-title",
			},
			expected: &librarydb.Book{Author: "some-author", Title: "some-title"},
			bookSvcFn: func(mockRepo *librarydb_mock.MockBookRepo) {
				mockRepo.EXPECT().
					Create(gomock.Any(), &librarydb.Book{Author: "some-author", Title: "some-title"}).
					Return(int64(1), nil)
				mockRepo.EXPECT().
					Find(gomock.Any(), dbkit.Equal("id", int64(1))).
					Return([]*librarydb.Book{{Author: "some-author", Title: "some-title"}}, nil)
			},
		},
	}

	for _, tt := range testcases {
		t.Run(tt.testName, func(t *testing.T) {
			svc, mock := createBookSvc(t, tt.bookSvcFn)
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

func TestBookSvc_RetrieveOne(t *testing.T) {
	testcases := []struct {
		testName    string
		bookSvcFn   bookSvcFn
		paramID     string
		expected    *librarydb.Book
		expectedErr string
	}{
		{
			paramID:     "",
			expectedErr: "paramID is missing",
		},
		{
			paramID: "1",
			bookSvcFn: func(mockRepo *librarydb_mock.MockBookRepo) {
				mockRepo.EXPECT().
					Find(gomock.Any(), dbkit.Equal("id", int64(1))).
					Return(nil, errors.New("some-error"))
			},
			expectedErr: "some-error",
		},
		{
			paramID: "1",
			bookSvcFn: func(mockRepo *librarydb_mock.MockBookRepo) {
				mockRepo.EXPECT().
					Find(gomock.Any(), dbkit.Equal("id", int64(1))).
					Return([]*librarydb.Book{
						{
							ID:    1,
							Title: "some-title",
						},
					}, nil)
			},
			expected: &librarydb.Book{
				ID:    1,
				Title: "some-title",
			},
		},
	}

	for _, tt := range testcases {
		t.Run(tt.testName, func(t *testing.T) {
			svc, mock := createBookSvc(t, tt.bookSvcFn)
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

func TestBookSvc_Retrieve(t *testing.T) {
	testcases := []struct {
		testName    string
		bookSvcFn   bookSvcFn
		expected    []*librarydb.Book
		expectedErr string
	}{}
	for _, tt := range testcases {
		t.Run(tt.testName, func(t *testing.T) {
			svc, mock := createBookSvc(t, tt.bookSvcFn)
			defer mock.Finish()

			books, err := svc.Find(context.Background())
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
		testName    string
		bookSvcFn   bookSvcFn
		paramID     string
		expectedErr string
	}{
		{
			paramID:     "",
			expectedErr: `paramID is missing`,
		},
		{
			paramID:     "1",
			expectedErr: `some-error`,
			bookSvcFn: func(mockRepo *librarydb_mock.MockBookRepo) {
				mockRepo.EXPECT().
					Delete(gomock.Any(), dbkit.Equal(librarydb.BookTable.ID, int64(1))).
					Return(int64(0), errors.New("some-error"))
			},
		},
		{
			paramID: "1",
			bookSvcFn: func(mockRepo *librarydb_mock.MockBookRepo) {
				mockRepo.EXPECT().
					Delete(gomock.Any(), dbkit.Equal(librarydb.BookTable.ID, int64(1))).
					Return(int64(1), nil)
			},
		},
		{
			testName: "success even if no affected row (idempotent)",
			paramID:  "1",
			bookSvcFn: func(mockRepo *librarydb_mock.MockBookRepo) {
				mockRepo.EXPECT().
					Delete(gomock.Any(), dbkit.Equal(librarydb.BookTable.ID, int64(1))).
					Return(int64(0), nil)
			},
		},
	}
	for _, tt := range testcases {
		t.Run(tt.testName, func(t *testing.T) {
			svc, mock := createBookSvc(t, tt.bookSvcFn)
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
		testName    string
		bookSvcFn   bookSvcFn
		paramID     string
		book        *librarydb.Book
		expected    *librarydb.Book
		expectedErr string
	}{
		{
			testName:    "empty paramID",
			paramID:     "",
			expectedErr: `paramID is missing`,
		},
		{
			testName:    "zero paramID",
			paramID:     "0",
			expectedErr: `paramID is missing`,
		},
		{
			testName:    "bad request",
			paramID:     "1",
			book:        &librarydb.Book{},
			expectedErr: "Key: 'Book.Title' Error:Field validation for 'Title' failed on the 'required' tag\nKey: 'Book.Author' Error:Field validation for 'Author' failed on the 'required' tag",
		},
		{
			testName:    "update error",
			paramID:     "1",
			book:        &librarydb.Book{Author: "some-author", Title: "some-title"},
			expectedErr: "update error",
			bookSvcFn: func(mockRepo *librarydb_mock.MockBookRepo) {
				mockRepo.EXPECT().
					Update(gomock.Any(), &librarydb.Book{Author: "some-author", Title: "some-title"}, dbkit.Equal("id", int64(1))).
					Return(int64(-1), errors.New("update error"))
			},
		},
		{
			testName:    "nothing to update",
			paramID:     "1",
			book:        &librarydb.Book{Author: "some-author", Title: "some-title"},
			expectedErr: "sql: no rows in result set",
			bookSvcFn: func(mockRepo *librarydb_mock.MockBookRepo) {
				mockRepo.EXPECT().
					Update(gomock.Any(), &librarydb.Book{Author: "some-author", Title: "some-title"}, dbkit.Equal("id", int64(1))).
					Return(int64(0), nil)
			},
		},
		{
			testName:    "Find error",
			paramID:     "1",
			book:        &librarydb.Book{Author: "some-author", Title: "some-title"},
			expectedErr: "Find-error",
			bookSvcFn: func(mockRepo *librarydb_mock.MockBookRepo) {
				mockRepo.EXPECT().
					Update(gomock.Any(), &librarydb.Book{Author: "some-author", Title: "some-title"}, dbkit.Equal("id", int64(1))).
					Return(int64(1), nil)
				mockRepo.EXPECT().
					Find(gomock.Any(), dbkit.Equal("id", int64(1))).
					Return(nil, errors.New("Find-error"))
			},
		},
		{
			testName:    "Find error",
			paramID:     "1",
			book:        &librarydb.Book{Author: "some-author", Title: "some-title"},
			expectedErr: "Find-error",
			bookSvcFn: func(mockRepo *librarydb_mock.MockBookRepo) {
				mockRepo.EXPECT().
					Update(gomock.Any(), &librarydb.Book{Author: "some-author", Title: "some-title"}, dbkit.Equal("id", int64(1))).
					Return(int64(1), nil)
				mockRepo.EXPECT().
					Find(gomock.Any(), dbkit.Equal("id", int64(1))).
					Return(nil, errors.New("Find-error"))
			},
		},
	}
	for _, tt := range testcases {
		t.Run(tt.testName, func(t *testing.T) {
			svc, mock := createBookSvc(t, tt.bookSvcFn)
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
		testName    string
		bookSvcFn   bookSvcFn
		paramID     string
		book        *librarydb.Book
		expected    *librarydb.Book
		expectedErr string
	}{
		{
			testName:    "empty paramID",
			paramID:     "",
			expectedErr: "paramID is missing",
		},
		{
			testName:    "zero paramID",
			paramID:     "0",
			expectedErr: "paramID is missing",
		},
		{
			testName:    "patch error",
			paramID:     "1",
			book:        &librarydb.Book{Author: "some-author", Title: "some-title"},
			expectedErr: "patch-error",
			bookSvcFn: func(mockRepo *librarydb_mock.MockBookRepo) {
				mockRepo.EXPECT().
					Patch(gomock.Any(), &librarydb.Book{Author: "some-author", Title: "some-title"}, dbkit.Equal("id", int64(1))).
					Return(int64(-1), errors.New("patch-error"))
			},
		},
		{
			testName:    "patch error",
			paramID:     "1",
			book:        &librarydb.Book{Author: "some-author", Title: "some-title"},
			expectedErr: "sql: no rows in result set",
			bookSvcFn: func(mockRepo *librarydb_mock.MockBookRepo) {
				mockRepo.EXPECT().
					Patch(gomock.Any(), &librarydb.Book{Author: "some-author", Title: "some-title"}, dbkit.Equal("id", int64(1))).
					Return(int64(0), nil)
			},
		},
		{
			testName:    "Find error",
			paramID:     "1",
			book:        &librarydb.Book{Author: "some-author", Title: "some-title"},
			expectedErr: "Find-error",
			bookSvcFn: func(mockRepo *librarydb_mock.MockBookRepo) {
				mockRepo.EXPECT().
					Patch(gomock.Any(), &librarydb.Book{Author: "some-author", Title: "some-title"}, dbkit.Equal("id", int64(1))).
					Return(int64(1), nil)
				mockRepo.EXPECT().
					Find(gomock.Any(), dbkit.Equal("id", int64(1))).
					Return(nil, errors.New("Find-error"))
			},
		},
		{
			paramID:  "1",
			book:     &librarydb.Book{Author: "some-author", Title: "some-title"},
			expected: &librarydb.Book{Author: "some-author", Title: "some-title"},
			bookSvcFn: func(mockRepo *librarydb_mock.MockBookRepo) {
				mockRepo.EXPECT().
					Patch(gomock.Any(), &librarydb.Book{Author: "some-author", Title: "some-title"}, dbkit.Equal("id", int64(1))).
					Return(int64(1), nil)
				mockRepo.EXPECT().
					Find(gomock.Any(), dbkit.Equal("id", int64(1))).
					Return([]*librarydb.Book{{Author: "some-author", Title: "some-title"}}, nil)
			},
		},
	}
	for _, tt := range testcases {
		t.Run(tt.testName, func(t *testing.T) {
			svc, mock := createBookSvc(t, tt.bookSvcFn)
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
