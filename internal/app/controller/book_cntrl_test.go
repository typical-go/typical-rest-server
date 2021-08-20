package controller_test

import (
	"errors"
	"fmt"
	"net/http"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/require"
	"github.com/typical-go/typical-rest-server/internal/app/controller"
	"github.com/typical-go/typical-rest-server/internal/app/repo"
	"github.com/typical-go/typical-rest-server/internal/app/service"
	"github.com/typical-go/typical-rest-server/internal/generated/mock/app/service_mock"
	"github.com/typical-go/typical-rest-server/pkg/echokit"
	"github.com/typical-go/typical-rest-server/pkg/echotest"
)

type (
	OnMockBookSvc     func(*service_mock.MockBookSvc)
	BookCntrlTestCase struct {
		TestName string
		echotest.TestCase
		OnMockBookSvc
	}
)

func CreateBookCntrl(t *testing.T, onMockBookSvc OnMockBookSvc) (*controller.BookCntrl, *gomock.Controller) {
	mock := gomock.NewController(t)
	mockSvc := service_mock.NewMockBookSvc(mock)
	if onMockBookSvc != nil {
		onMockBookSvc(mockSvc)
	}
	return &controller.BookCntrl{Svc: mockSvc}, mock
}

func TestBookCntrl_SetRoute(t *testing.T) {
	e := echo.New()
	echokit.SetRoute(e, &controller.BookCntrl{})
	require.Equal(t, []string{
		"/books\tGET,POST",
		"/books/:id\tDELETE,GET,HEAD,PATCH,PUT",
	}, echokit.DumpEcho(e))
}

func TestBookCntrl_FindOne(t *testing.T) {
	testcases := []BookCntrlTestCase{
		{
			TestName: "valid ID",
			TestCase: echotest.TestCase{
				Request: echotest.Request{
					Method:    http.MethodGet,
					Target:    "/",
					URLParams: map[string]string{"id": "1"},
				},
				ExpectedResponse: echotest.Response{
					Code: http.StatusOK,
					Body: "{\"id\":1,\"title\":\"title1\",\"author\":\"author1\",\"update_at\":\"0001-01-01T00:00:00Z\",\"created_at\":\"0001-01-01T00:00:00Z\"}\n",
					Header: http.Header{
						"Content-Type": {"application/json; charset=UTF-8"},
					},
				},
			},
			OnMockBookSvc: func(svc *service_mock.MockBookSvc) {
				svc.EXPECT().FindOne(gomock.Any(), "1").Return(&repo.Book{ID: 1, Title: "title1", Author: "author1"}, nil)
			},
		},
		{
			TestName: "repo not found",
			TestCase: echotest.TestCase{
				Request: echotest.Request{
					Method:    http.MethodGet,
					Target:    "/",
					URLParams: map[string]string{"id": "3"},
				},
				ExpectedError: "code=404, message=Not Found",
			},
			OnMockBookSvc: func(svc *service_mock.MockBookSvc) {
				svc.EXPECT().FindOne(gomock.Any(), "3").Return(nil, echo.NewHTTPError(404))
			},
		},
		{
			TestName: "validation error",
			TestCase: echotest.TestCase{
				Request: echotest.Request{
					Method:    http.MethodGet,
					Target:    "/",
					URLParams: map[string]string{"id": "2"},
				},
				ExpectedError: "code=422, message=some-validation",
			},
			OnMockBookSvc: func(svc *service_mock.MockBookSvc) {
				svc.EXPECT().FindOne(gomock.Any(), "2").Return(nil, echokit.NewValidErr("some-validation"))
			},
		},
		{
			TestName: "error",
			TestCase: echotest.TestCase{
				Request: echotest.Request{
					Method:    http.MethodGet,
					Target:    "/",
					URLParams: map[string]string{"id": "2"},
				},
				ExpectedError: "code=500, message=some-error",
			},
			OnMockBookSvc: func(svc *service_mock.MockBookSvc) {
				svc.EXPECT().FindOne(gomock.Any(), "2").Return(nil, errors.New("some-error"))
			},
		},
	}

	for _, tt := range testcases {
		t.Run(tt.TestName, func(t *testing.T) {
			cntrl, mock := CreateBookCntrl(t, tt.OnMockBookSvc)
			defer mock.Finish()
			tt.Execute(t, cntrl.FindOne)
		})
	}
}

func TestBookCntrl_Find(t *testing.T) {
	testcases := []BookCntrlTestCase{
		{
			TestCase: echotest.TestCase{
				Request: echotest.Request{
					Method: http.MethodGet,
					Target: "/",
				},
				ExpectedResponse: echotest.Response{
					Code: http.StatusOK,
					Body: "[{\"id\":1,\"title\":\"title1\",\"author\":\"author1\",\"update_at\":\"0001-01-01T00:00:00Z\",\"created_at\":\"0001-01-01T00:00:00Z\"},{\"id\":2,\"title\":\"title2\",\"author\":\"author2\",\"update_at\":\"0001-01-01T00:00:00Z\",\"created_at\":\"0001-01-01T00:00:00Z\"}]\n",
					Header: http.Header{
						"Content-Type":  {"application/json; charset=UTF-8"},
						"X-Total-Count": {"10"},
					},
				},
			},
			OnMockBookSvc: func(svc *service_mock.MockBookSvc) {
				svc.EXPECT().
					Find(gomock.Any(), &service.FindBookReq{}).
					Return(&service.FindBookResp{
						TotalCount: "10",
						Books: []*repo.Book{
							{ID: 1, Title: "title1", Author: "author1"},
							{ID: 2, Title: "title2", Author: "author2"},
						},
					}, nil)
			},
		},
		{
			TestCase: echotest.TestCase{
				Request: echotest.Request{
					Method: http.MethodGet,
					Target: "/?limit=20&offset=10",
				},
				ExpectedError: "code=500, message=some-error",
			},
			OnMockBookSvc: func(svc *service_mock.MockBookSvc) {
				svc.EXPECT().
					Find(gomock.Any(), &service.FindBookReq{Limit: 20, Offset: 10}).
					Return(nil, fmt.Errorf("some-error"))
			},
		},
		{
			TestCase: echotest.TestCase{
				Request: echotest.Request{
					Method: http.MethodGet,
					Target: "/?sort=name,created_at",
				},
				ExpectedError: "code=500, message=some-error",
			},
			OnMockBookSvc: func(svc *service_mock.MockBookSvc) {
				svc.EXPECT().
					Find(gomock.Any(), &service.FindBookReq{Sort: "name,created_at"}).
					Return(nil, fmt.Errorf("some-error"))
			},
		},
	}

	for _, tt := range testcases {
		t.Run(tt.TestName, func(t *testing.T) {
			cntrl, mock := CreateBookCntrl(t, tt.OnMockBookSvc)
			defer mock.Finish()
			tt.Execute(t, cntrl.Find)
		})
	}
}

func TestBookCntrl_Update(t *testing.T) {
	testcases := []BookCntrlTestCase{
		{
			TestCase: echotest.TestCase{
				Request: echotest.Request{
					Method:    http.MethodPut,
					Target:    "/",
					URLParams: map[string]string{"id": "1"},
					Header:    http.Header{"Content-Type": {"application/json"}},
					Body:      `{bad-json`,
				},
				ExpectedError: "code=400, message=Syntax error: offset=2, error=invalid character 'b' looking for beginning of object key string, internal=invalid character 'b' looking for beginning of object key string",
			},
		},
		{
			TestCase: echotest.TestCase{
				Request: echotest.Request{
					Method:    http.MethodPut,
					Target:    "/",
					URLParams: map[string]string{"id": "1"},
					Header:    http.Header{"Content-Type": {"application/json"}},
					Body:      `{"title":"some-title", "author": "some-author"}`,
				},
				ExpectedError: "code=500, message=some-error",
			},
			OnMockBookSvc: func(svc *service_mock.MockBookSvc) {
				svc.EXPECT().
					Update(gomock.Any(), "1", &repo.Book{ID: 1, Title: "some-title", Author: "some-author"}).
					Return(nil, errors.New("some-error"))
			},
		},
		{
			TestCase: echotest.TestCase{
				Request: echotest.Request{
					Method:    http.MethodPut,
					Target:    "/",
					URLParams: map[string]string{"id": "1"},
					Header:    http.Header{"Content-Type": {"application/json"}},
					Body:      `{"title":"some-title", "author": "some-author"}`,
				},
				ExpectedResponse: echotest.Response{
					Code: http.StatusOK,
					Body: "{\"id\":1,\"title\":\"some-title\",\"author\":\"some-author\",\"update_at\":\"0001-01-01T00:00:00Z\",\"created_at\":\"0001-01-01T00:00:00Z\"}\n",
					Header: http.Header{
						"Content-Type": {"application/json; charset=UTF-8"},
					},
				},
			},
			OnMockBookSvc: func(svc *service_mock.MockBookSvc) {
				svc.EXPECT().
					Update(gomock.Any(), "1", &repo.Book{ID: 1, Title: "some-title", Author: "some-author"}).
					Return(&repo.Book{ID: 1, Title: "some-title", Author: "some-author"}, nil)
			},
		},
	}

	for _, tt := range testcases {
		t.Run(tt.TestName, func(t *testing.T) {
			cntrl, mock := CreateBookCntrl(t, tt.OnMockBookSvc)
			defer mock.Finish()
			tt.Execute(t, cntrl.Update)
		})
	}
}

func TestBookCntrl_Patch(t *testing.T) {
	testcases := []BookCntrlTestCase{
		{
			TestCase: echotest.TestCase{
				Request: echotest.Request{
					Method:    http.MethodPut,
					Target:    "/",
					URLParams: map[string]string{"id": "1"},
					Header:    http.Header{"Content-Type": {"application/json"}},
					Body:      `{bad-json`,
				},
				ExpectedError: "code=400, message=Syntax error: offset=2, error=invalid character 'b' looking for beginning of object key string, internal=invalid character 'b' looking for beginning of object key string",
			},
		},
		{
			TestCase: echotest.TestCase{
				Request: echotest.Request{
					Method:    http.MethodPut,
					Target:    "/",
					URLParams: map[string]string{"id": "1"},
					Header:    http.Header{"Content-Type": {"application/json"}},
					Body:      `{"title":"some-title", "author": "some-author"}`,
				},
				ExpectedError: "code=500, message=some-error",
			},
			OnMockBookSvc: func(svc *service_mock.MockBookSvc) {
				svc.EXPECT().
					Patch(gomock.Any(), "1", &repo.Book{ID: 1, Title: "some-title", Author: "some-author"}).
					Return(nil, errors.New("some-error"))
			},
		},
		{
			TestCase: echotest.TestCase{
				Request: echotest.Request{
					Method:    http.MethodPut,
					Target:    "/",
					URLParams: map[string]string{"id": "1"},
					Header:    http.Header{"Content-Type": {"application/json"}},
					Body:      `{"title":"some-title", "author": "some-author"}`,
				},
				ExpectedResponse: echotest.Response{
					Code: http.StatusOK,
					Body: "{\"id\":1,\"title\":\"some-title\",\"author\":\"some-author\",\"update_at\":\"0001-01-01T00:00:00Z\",\"created_at\":\"0001-01-01T00:00:00Z\"}\n",
					Header: http.Header{
						"Content-Type": {"application/json; charset=UTF-8"},
					},
				},
			},
			OnMockBookSvc: func(svc *service_mock.MockBookSvc) {
				svc.EXPECT().
					Patch(gomock.Any(), "1", &repo.Book{ID: 1, Title: "some-title", Author: "some-author"}).
					Return(&repo.Book{ID: 1, Title: "some-title", Author: "some-author"}, nil)
			},
		},
	}

	for _, tt := range testcases {
		t.Run(tt.TestName, func(t *testing.T) {
			cntrl, mock := CreateBookCntrl(t, tt.OnMockBookSvc)
			defer mock.Finish()
			tt.Execute(t, cntrl.Patch)
		})
	}
}

func TestBookCntrl_Delete(t *testing.T) {
	testcases := []BookCntrlTestCase{
		{
			TestCase: echotest.TestCase{
				Request: echotest.Request{
					Method:    http.MethodDelete,
					Target:    "/",
					URLParams: map[string]string{"id": "1"},
				},
				ExpectedResponse: echotest.Response{
					Code:   http.StatusNoContent,
					Header: http.Header{},
				},
			},
			OnMockBookSvc: func(svc *service_mock.MockBookSvc) {
				svc.EXPECT().Delete(gomock.Any(), "1").Return(nil)
			},
		},
		{
			TestCase: echotest.TestCase{
				Request: echotest.Request{
					Method:    http.MethodDelete,
					Target:    "/",
					URLParams: map[string]string{"id": "1"},
				},
				ExpectedError: "code=500, message=some-error",
			},
			OnMockBookSvc: func(svc *service_mock.MockBookSvc) {
				svc.EXPECT().Delete(gomock.Any(), "1").Return(errors.New("some-error"))
			},
		},
		{
			TestCase: echotest.TestCase{
				Request: echotest.Request{
					Method:    http.MethodDelete,
					Target:    "/",
					URLParams: map[string]string{"id": "1"},
				},
				ExpectedError: "code=422, message=some-validation",
			},
			OnMockBookSvc: func(svc *service_mock.MockBookSvc) {
				svc.EXPECT().Delete(gomock.Any(), "1").Return(echokit.NewValidErr("some-validation"))
			},
		},
	}

	for _, tt := range testcases {
		t.Run(tt.TestName, func(t *testing.T) {
			cntrl, mock := CreateBookCntrl(t, tt.OnMockBookSvc)
			defer mock.Finish()
			tt.Execute(t, cntrl.Delete)
		})
	}
}

func TestBookCntrl_Create(t *testing.T) {
	testcases := []BookCntrlTestCase{
		{
			TestCase: echotest.TestCase{
				Request: echotest.Request{
					Method: http.MethodPost,
					Target: "/",
					Body:   `invalid}`,
					Header: http.Header{"Content-Type": {"application/json"}},
				},
				ExpectedError: "code=400, message=Syntax error: offset=1, error=invalid character 'i' looking for beginning of value, internal=invalid character 'i' looking for beginning of value",
			},
		},
		{
			TestCase: echotest.TestCase{
				Request: echotest.Request{
					Method: http.MethodPost,
					Target: "/",
					Body:   `{"author":"some-author", "title":"some-title"}`,
					Header: http.Header{"Content-Type": {"application/json"}},
				},
				ExpectedError: "code=500, message=some-error",
			},
			OnMockBookSvc: func(svc *service_mock.MockBookSvc) {
				svc.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf("some-error"))
			},
		},
		{
			TestCase: echotest.TestCase{
				Request: echotest.Request{
					Method: http.MethodPost,
					Target: "/",
					Body:   `{"author":"some-author", "title":"some-title"}`,
					Header: http.Header{"Content-Type": {"application/json"}},
				},
				ExpectedResponse: echotest.Response{
					Body: "{\"id\":999,\"title\":\"some-title\",\"author\":\"some-author\",\"update_at\":\"0001-01-01T00:00:00Z\",\"created_at\":\"0001-01-01T00:00:00Z\"}\n",
					Code: http.StatusCreated,
					Header: http.Header{
						"Content-Type": {"application/json; charset=UTF-8"},
						"Location":     {"/books/999"},
					},
				},
			},
			OnMockBookSvc: func(svc *service_mock.MockBookSvc) {
				svc.EXPECT().
					Create(gomock.Any(), gomock.Any()).
					Return(&repo.Book{ID: 999, Author: "some-author", Title: "some-title"}, nil)
			},
		},
	}

	for _, tt := range testcases {
		t.Run(tt.TestName, func(t *testing.T) {
			cntrl, mock := CreateBookCntrl(t, tt.OnMockBookSvc)
			defer mock.Finish()
			tt.Execute(t, cntrl.Create)
		})
	}
}
