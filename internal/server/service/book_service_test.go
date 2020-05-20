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
		testName    string
		builder     bookSvcBuilder
		paramID     string
		book        *repository.Book
		expectedErr string
	}
)

func (b *bookSvcBuilder) build(mock *gomock.Controller) *service.BookServiceImpl {
	mockRepo := repository_mock.NewMockBookRepo(mock)
	if b.mockRepoFn != nil {
		b.mockRepoFn(mockRepo)
	}

	return &service.BookServiceImpl{
		BookRepo: mockRepo,
	}
}

func TestBookService_FindOne(t *testing.T) {
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
			mock := gomock.NewController(t)
			defer mock.Finish()

			book, err := tt.build(mock).FindOne(context.Background(), tt.paramID)
			if tt.expectedErr != "" {
				require.EqualError(t, err, tt.expectedErr)
			} else {
				require.NoError(t, err)
			}
			require.Equal(t, tt.expected, book)
		})
	}
}

func TestBookService_Find(t *testing.T) {
	testcases := []findTestCase{}
	for _, tt := range testcases {
		t.Run(tt.testName, func(t *testing.T) {
			mock := gomock.NewController(t)
			defer mock.Finish()

			books, err := tt.build(mock).Find(context.Background())

			if tt.expectedErr != "" {
				require.EqualError(t, err, tt.expectedErr)
			} else {
				require.NoError(t, err)
			}
			require.Equal(t, tt.expected, books)
		})
	}
}

func TestBookService_Create(t *testing.T) {
	testcases := []createTestCase{}

	for _, tt := range testcases {
		t.Run(tt.testName, func(t *testing.T) {
			mock := gomock.NewController(t)
			defer mock.Finish()

			id, err := tt.build(mock).Create(context.Background(), tt.book)

			if tt.expectedErr != "" {
				require.EqualError(t, err, tt.expectedErr)
			} else {
				require.NoError(t, err)
			}
			require.Equal(t, tt.expected, id)
		})
	}
}

func TestBookService_Delete(t *testing.T) {
	testcases := []deleteTestCase{
		{
			paramID:     "",
			expectedErr: `Validation: strconv.ParseInt: parsing "": invalid syntax`,
		},
	}
	for _, tt := range testcases {
		t.Run(tt.testName, func(t *testing.T) {
			mock := gomock.NewController(t)
			defer mock.Finish()

			err := tt.build(mock).Delete(context.Background(), tt.paramID)

			if tt.expectedErr != "" {
				require.EqualError(t, err, tt.expectedErr)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
func TestBookService_Update(t *testing.T) {
	testcases := []updateTestCase{
		{
			paramID:     "",
			expectedErr: `Validation: strconv.ParseInt: parsing "": invalid syntax`,
		},
	}
	for _, tt := range testcases {
		t.Run(tt.testName, func(t *testing.T) {
			mock := gomock.NewController(t)
			defer mock.Finish()

			svc := tt.builder.build(mock)
			err := svc.Update(context.Background(), tt.paramID, tt.book)

			if tt.expectedErr != "" {
				require.EqualError(t, err, tt.expectedErr)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
