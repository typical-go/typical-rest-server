package service_test

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"github.com/typical-go/typical-rest-server/internal/server/repository"
	"github.com/typical-go/typical-rest-server/internal/server/repository_mock"
	"github.com/typical-go/typical-rest-server/internal/server/service"
	"github.com/typical-go/typical-rest-server/pkg/dbkit"
)

type (
	bookSvcBuilder struct {
		mockRepoFn func(mockRepo *repository_mock.MockBookRepo)
	}

	findOneTestCase struct {
		testName string
		bookSvcBuilder
		paramID     string
		expected    *repository.Book
		expectedErr string
	}

	findTestCase struct {
		testName string
		bookSvcBuilder
		expected    []*repository.Book
		expectedErr string
	}

	createTestCase struct {
		testName string
		bookSvcBuilder
		book        *repository.Book
		expected    int64
		expectedErr string
	}

	deleteTestCase struct {
		testName string
		bookSvcBuilder
		paramID     string
		expectedErr string
	}

	updateTestCase struct {
		testName string
		bookSvcBuilder
		paramID     string
		book        *repository.Book
		expectedErr string
	}
)

func (b *bookSvcBuilder) build(t *testing.T) (*service.BookSvcImpl, *gomock.Controller) {
	mock := gomock.NewController(t)
	mockRepo := repository_mock.NewMockBookRepo(mock)
	if b.mockRepoFn != nil {
		b.mockRepoFn(mockRepo)
	}

	return &service.BookSvcImpl{
		BookRepo: mockRepo,
	}, mock
}

func TestBookSvc_FindOne(t *testing.T) {
	testcases := []findOneTestCase{
		{
			paramID:     "",
			expectedErr: `Validation: strconv.ParseInt: parsing "": invalid syntax`,
		},
		{
			paramID: "1",
			bookSvcBuilder: bookSvcBuilder{
				mockRepoFn: func(mockRepo *repository_mock.MockBookRepo) {
					mockRepo.EXPECT().
						Find(gomock.Any(), dbkit.Equal("id", int64(1))).
						Return(nil, errors.New("some-error"))
				},
			},
			expectedErr: "some-error",
		},
		{
			paramID: "1",
			bookSvcBuilder: bookSvcBuilder{
				mockRepoFn: func(mockRepo *repository_mock.MockBookRepo) {
					mockRepo.EXPECT().
						Find(gomock.Any(), dbkit.Equal("id", int64(1))).
						Return([]*repository.Book{
							{
								ID:    1,
								Title: "some-title",
							},
						}, nil)
				},
			},
			expected: &repository.Book{
				ID:    1,
				Title: "some-title",
			},
		},
	}

	for _, tt := range testcases {
		t.Run(tt.testName, func(t *testing.T) {
			svc, mock := tt.build(t)
			defer mock.Finish()

			book, err := svc.FindOne(context.Background(), tt.paramID)
			if tt.expectedErr != "" {
				require.EqualError(t, err, tt.expectedErr)
				return
			}

			require.NoError(t, err)
			require.Equal(t, tt.expected, book)
		})
	}
}

func TestBookSvc_Find(t *testing.T) {
	testcases := []findTestCase{}
	for _, tt := range testcases {
		t.Run(tt.testName, func(t *testing.T) {
			svc, mock := tt.build(t)
			defer mock.Finish()

			books, err := svc.Find(context.Background())
			if tt.expectedErr != "" {
				require.EqualError(t, err, tt.expectedErr)
				return
			}

			require.NoError(t, err)
			require.Equal(t, tt.expected, books)
		})
	}
}

func TestBookSvc_Create(t *testing.T) {
	testcases := []createTestCase{}

	for _, tt := range testcases {
		t.Run(tt.testName, func(t *testing.T) {
			svc, mock := tt.build(t)
			defer mock.Finish()

			id, err := svc.Create(context.Background(), tt.book)
			if tt.expectedErr != "" {
				require.EqualError(t, err, tt.expectedErr)
				return
			}

			require.NoError(t, err)
			require.Equal(t, tt.expected, id)
		})
	}
}

func TestBookSvc_Delete(t *testing.T) {
	testcases := []deleteTestCase{
		{
			paramID:     "",
			expectedErr: `Validation: paramID is missing`,
		},
	}
	for _, tt := range testcases {
		t.Run(tt.testName, func(t *testing.T) {
			svc, mock := tt.build(t)
			defer mock.Finish()

			err := svc.Delete(context.Background(), tt.paramID)
			if tt.expectedErr != "" {
				require.EqualError(t, err, tt.expectedErr)
				return
			}
			require.NoError(t, err)
		})
	}
}
func TestBookSvc_Update(t *testing.T) {
	testcases := []updateTestCase{
		{
			paramID:     "",
			expectedErr: `Validation: paramID is missing`,
		},
	}
	for _, tt := range testcases {
		t.Run(tt.testName, func(t *testing.T) {
			svc, mock := tt.build(t)
			defer mock.Finish()

			err := svc.Update(context.Background(), tt.paramID, tt.book)
			if tt.expectedErr != "" {
				require.EqualError(t, err, tt.expectedErr)
				return
			}
			require.NoError(t, err)
		})
	}
}
