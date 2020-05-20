package controller_test

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/typical-go/typical-rest-server/internal/server/controller"
	"github.com/typical-go/typical-rest-server/internal/server/repository"
	"github.com/typical-go/typical-rest-server/internal/server/service_mock"
	"github.com/typical-go/typical-rest-server/pkg/echotest"
	"github.com/typical-go/typical-rest-server/pkg/errvalid"
)

type (
	testCase struct {
		testName string
		echotest.TestCase
		bookCntrlBuilder
	}

	bookCntrlBuilder struct {
		bookSvcFn func(*service_mock.MockBookSvc)
	}
)

func (c *bookCntrlBuilder) build(mock *gomock.Controller) *controller.BookCntrl {
	mockSvc := service_mock.NewMockBookSvc(mock)
	if c.bookSvcFn != nil {
		c.bookSvcFn(mockSvc)
	}
	return &controller.BookCntrl{
		BookSvc: mockSvc,
	}
}

func TestBookController_FindOne(t *testing.T) {
	testcases := []testCase{
		{
			testName: "valid ID",
			TestCase: echotest.TestCase{
				Request: echotest.Request{
					Method:    http.MethodGet,
					Target:    "/",
					URLParams: map[string]string{"id": "1"},
				},
				ExpectedCode: http.StatusOK,
				ExpectedBody: "{\"id\":1,\"title\":\"title1\",\"author\":\"author1\",\"update_at\":\"0001-01-01T00:00:00Z\",\"created_at\":\"0001-01-01T00:00:00Z\"}\n",
			},
			bookCntrlBuilder: bookCntrlBuilder{
				bookSvcFn: func(svc *service_mock.MockBookSvc) {
					svc.EXPECT().FindOne(gomock.Any(), "1").Return(&repository.Book{ID: 1, Title: "title1", Author: "author1"}, nil)
				},
			},
		},
		{
			testName: "entity not found",
			TestCase: echotest.TestCase{
				Request: echotest.Request{
					Method:    http.MethodGet,
					Target:    "/",
					URLParams: map[string]string{"id": "3"},
				},
				ExpectedErr: "code=404, message=Not Found",
			},
			bookCntrlBuilder: bookCntrlBuilder{
				bookSvcFn: func(svc *service_mock.MockBookSvc) {
					svc.EXPECT().FindOne(gomock.Any(), "3").Return(nil, sql.ErrNoRows)
				},
			},
		},
		{
			testName: "validation error",
			TestCase: echotest.TestCase{
				Request: echotest.Request{
					Method:    http.MethodGet,
					Target:    "/",
					URLParams: map[string]string{"id": "2"},
				},
				ExpectedErr: "code=422, message=some-validation",
			},
			bookCntrlBuilder: bookCntrlBuilder{
				bookSvcFn: func(svc *service_mock.MockBookSvc) {
					svc.EXPECT().FindOne(gomock.Any(), "2").Return(nil, errvalid.New("some-validation"))
				},
			},
		},
		{
			testName: "error",
			TestCase: echotest.TestCase{
				Request: echotest.Request{
					Method:    http.MethodGet,
					Target:    "/",
					URLParams: map[string]string{"id": "2"},
				},
				ExpectedErr: "code=500, message=some-error",
			},
			bookCntrlBuilder: bookCntrlBuilder{
				bookSvcFn: func(svc *service_mock.MockBookSvc) {
					svc.EXPECT().FindOne(gomock.Any(), "2").Return(nil, errors.New("some-error"))
				},
			},
		},
	}

	for _, tt := range testcases {
		t.Run(tt.testName, func(t *testing.T) {
			mock := gomock.NewController(t)
			defer mock.Finish()
			tt.Execute(t, tt.build(mock).FindOne)
		})
	}
}

func TestBookController_Find(t *testing.T) {
	testcases := []testCase{
		{
			TestCase: echotest.TestCase{
				Request: echotest.Request{
					Method: http.MethodGet,
					Target: "/",
				},
				ExpectedCode: http.StatusOK,
				ExpectedBody: "[{\"id\":1,\"title\":\"title1\",\"author\":\"author1\",\"update_at\":\"0001-01-01T00:00:00Z\",\"created_at\":\"0001-01-01T00:00:00Z\"},{\"id\":2,\"title\":\"title2\",\"author\":\"author2\",\"update_at\":\"0001-01-01T00:00:00Z\",\"created_at\":\"0001-01-01T00:00:00Z\"}]\n",
			},
			bookCntrlBuilder: bookCntrlBuilder{
				bookSvcFn: func(svc *service_mock.MockBookSvc) {
					svc.EXPECT().Find(gomock.Any()).Return([]*repository.Book{
						&repository.Book{ID: 1, Title: "title1", Author: "author1"},
						&repository.Book{ID: 2, Title: "title2", Author: "author2"},
					}, nil)
				},
			},
		},
		{
			TestCase: echotest.TestCase{
				Request: echotest.Request{
					Method: http.MethodGet,
					Target: "/",
				},
				ExpectedErr: "code=500, message=some-error",
			},
			bookCntrlBuilder: bookCntrlBuilder{
				bookSvcFn: func(svc *service_mock.MockBookSvc) {
					svc.EXPECT().Find(gomock.Any()).Return(nil, fmt.Errorf("some-error"))
				},
			},
		},
	}

	for _, tt := range testcases {
		t.Run(tt.testName, func(t *testing.T) {
			mock := gomock.NewController(t)
			defer mock.Finish()
			tt.Execute(t, tt.build(mock).Find)
		})
	}
}

func TestBookController_Update(t *testing.T) {
	testcases := []testCase{
		// TODO:
	}

	for _, tt := range testcases {
		t.Run(tt.testName, func(t *testing.T) {
			mock := gomock.NewController(t)
			defer mock.Finish()
			tt.Execute(t, tt.build(mock).Update)
		})
	}
}

func TestBookController_Delete(t *testing.T) {
	testcases := []testCase{
		{
			TestCase: echotest.TestCase{
				Request: echotest.Request{
					Method:    http.MethodDelete,
					Target:    "/",
					URLParams: map[string]string{"id": "1"},
				},
				ExpectedCode: http.StatusNoContent,
			},
			bookCntrlBuilder: bookCntrlBuilder{
				bookSvcFn: func(svc *service_mock.MockBookSvc) {
					svc.EXPECT().Delete(gomock.Any(), "1").Return(nil)
				},
			},
		},
		{
			TestCase: echotest.TestCase{
				Request: echotest.Request{
					Method:    http.MethodDelete,
					Target:    "/",
					URLParams: map[string]string{"id": "1"},
				},
				ExpectedErr: "code=500, message=some-error",
			},
			bookCntrlBuilder: bookCntrlBuilder{
				bookSvcFn: func(svc *service_mock.MockBookSvc) {
					svc.EXPECT().Delete(gomock.Any(), "1").Return(errors.New("some-error"))
				},
			},
		},
		{
			TestCase: echotest.TestCase{
				Request: echotest.Request{
					Method:    http.MethodDelete,
					Target:    "/",
					URLParams: map[string]string{"id": "1"},
				},
				ExpectedErr: "code=422, message=some-validation",
			},
			bookCntrlBuilder: bookCntrlBuilder{
				bookSvcFn: func(svc *service_mock.MockBookSvc) {
					svc.EXPECT().Delete(gomock.Any(), "1").Return(errvalid.New("some-validation"))
				},
			},
		},
	}

	for _, tt := range testcases {
		t.Run(tt.testName, func(t *testing.T) {
			mock := gomock.NewController(t)
			defer mock.Finish()
			tt.Execute(t, tt.build(mock).Delete)
		})
	}
}

func TestBookController_Create(t *testing.T) {
	testcases := []testCase{
		{
			TestCase: echotest.TestCase{
				Request: echotest.Request{
					Method: http.MethodPost,
					Target: "/",
					Body:   `invalid}`,
					Header: echotest.HeaderForJSON(),
				},
				ExpectedErr: `code=400, message=Syntax error: offset=1, error=invalid character 'i' looking for beginning of value`,
			},
		},
		{
			TestCase: echotest.TestCase{
				Request: echotest.Request{
					Method: http.MethodPost,
					Target: "/",
					Body:   `{"author":"some-author", "title":"some-title"}`,
					Header: echotest.HeaderForJSON(),
				},
				ExpectedErr: "code=500, message=some-error",
			},
			bookCntrlBuilder: bookCntrlBuilder{
				bookSvcFn: func(svc *service_mock.MockBookSvc) {
					svc.EXPECT().Create(gomock.Any(), gomock.Any()).Return(int64(-1), fmt.Errorf("some-error"))
				},
			},
		},
		{
			TestCase: echotest.TestCase{
				Request: echotest.Request{
					Method: http.MethodPost,
					Target: "/",
					Body:   `{"author":"some-author", "title":"some-title"}`,
					Header: echotest.HeaderForJSON(),
				},
				ExpectedCode: http.StatusCreated,
				ExpectedHeader: map[string]string{
					"Location": "/books/999",
				},
			},
			bookCntrlBuilder: bookCntrlBuilder{
				bookSvcFn: func(svc *service_mock.MockBookSvc) {
					svc.EXPECT().Create(gomock.Any(), gomock.Any()).Return(int64(999), nil)
				},
			},
		},
	}

	for _, tt := range testcases {
		t.Run(tt.testName, func(t *testing.T) {
			mock := gomock.NewController(t)
			defer mock.Finish()
			tt.Execute(t, tt.build(mock).Create)
		})
	}
}
