package service_test

import (
	"context"
	"errors"
	"testing"

	"github.com/typical-go/typical-rest-server/pkg/dbkit"

	"github.com/stretchr/testify/require"

	"github.com/golang/mock/gomock"
	"github.com/typical-go/typical-rest-server/server/repository_mock"
	"github.com/typical-go/typical-rest-server/server/service"

	"github.com/typical-go/typical-rest-server/server/repository"
)

var (
	findOneTestCases = []findOneTestCase{
		{
			paramID:     "",
			expectedErr: `Validation: strconv.ParseInt: parsing "": invalid syntax`,
		},
		{
			paramID: "1",
			bookSvcCtrl: bookSvcCtrl{
				mockRepoFn: func(mockRepo *repository_mock.MockBookRepo) {
					mockRepo.EXPECT().
						FindOne(gomock.Any(), int64(1)).
						Return(nil, errors.New("some-error"))
				},
			},
			expectedErr: "some-error",
		},
		{
			paramID: "1",
			bookSvcCtrl: bookSvcCtrl{
				mockRepoFn: func(mockRepo *repository_mock.MockBookRepo) {
					mockRepo.EXPECT().
						FindOne(gomock.Any(), int64(1)).
						Return(&repository.Book{
							ID:    1,
							Title: "some-title",
						}, nil)
				},
			},
			expected: &repository.Book{
				ID:    1,
				Title: "some-title",
			},
		},
	}

	findTestCases = []findTestCase{}

	createTestCases = []createTestCase{}

	deleteTestCases = []deleteTestCase{
		{
			paramID:     "",
			expectedErr: `Validation: strconv.ParseInt: parsing "": invalid syntax`,
		},
	}

	updateTestCases = []updateTestCase{
		{
			paramID:     "",
			expectedErr: `Validation: strconv.ParseInt: parsing "": invalid syntax`,
		},
	}
)

type (
	bookSvcCtrl struct {
		mockRepoFn func(mockRepo *repository_mock.MockBookRepo)
	}

	findOneTestCase struct {
		testName string
		bookSvcCtrl
		paramID     string
		expected    *repository.Book
		expectedErr string
	}

	findTestCase struct {
		testName string
		bookSvcCtrl
		opts        []dbkit.FindOption
		expected    []*repository.Book
		expectedErr string
	}

	createTestCase struct {
		testName string
		bookSvcCtrl
		book        *repository.Book
		expected    *repository.Book
		expectedErr string
	}

	deleteTestCase struct {
		testName string
		bookSvcCtrl
		paramID     string
		expectedErr string
	}

	updateTestCase struct {
		testName string
		bookSvcCtrl
		paramID     string
		book        *repository.Book
		expected    *repository.Book
		expectedErr string
	}
)

func createBookSvc(mock *gomock.Controller, ctrl bookSvcCtrl) *service.BookServiceImpl {
	mockRepo := repository_mock.NewMockBookRepo(mock)
	if ctrl.mockRepoFn != nil {
		ctrl.mockRepoFn(mockRepo)
	}

	return &service.BookServiceImpl{
		BookRepo: mockRepo,
	}
}

func TestBookService_FindOne(t *testing.T) {
	for _, tt := range findOneTestCases {
		t.Run(tt.testName, func(t *testing.T) {
			mock := gomock.NewController(t)
			defer mock.Finish()

			svc := createBookSvc(mock, tt.bookSvcCtrl)

			book, err := svc.FindOne(context.Background(), tt.paramID)
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
	for _, tt := range findTestCases {
		t.Run(tt.testName, func(t *testing.T) {
			mock := gomock.NewController(t)
			defer mock.Finish()

			svc := createBookSvc(mock, tt.bookSvcCtrl)
			books, err := svc.Find(context.Background(), tt.opts...)

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
	for _, tt := range createTestCases {
		t.Run(tt.testName, func(t *testing.T) {
			mock := gomock.NewController(t)
			defer mock.Finish()

			svc := createBookSvc(mock, tt.bookSvcCtrl)
			book, err := svc.Create(context.Background(), tt.book)

			if tt.expectedErr != "" {
				require.EqualError(t, err, tt.expectedErr)
			} else {
				require.NoError(t, err)
			}

			require.Equal(t, tt.expected, book)
		})
	}
}

func TestBookService_Delete(t *testing.T) {
	for _, tt := range deleteTestCases {
		t.Run(tt.testName, func(t *testing.T) {
			mock := gomock.NewController(t)
			defer mock.Finish()

			svc := createBookSvc(mock, tt.bookSvcCtrl)
			err := svc.Delete(context.Background(), tt.paramID)

			if tt.expectedErr != "" {
				require.EqualError(t, err, tt.expectedErr)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
func TestBookService_Update(t *testing.T) {
	for _, tt := range updateTestCases {
		t.Run(tt.testName, func(t *testing.T) {
			mock := gomock.NewController(t)
			defer mock.Finish()

			svc := createBookSvc(mock, tt.bookSvcCtrl)
			book, err := svc.Update(context.Background(), tt.paramID, tt.book)

			if tt.expectedErr != "" {
				require.EqualError(t, err, tt.expectedErr)
			} else {
				require.NoError(t, err)
			}

			require.Equal(t, tt.expected, book)
		})
	}
}
