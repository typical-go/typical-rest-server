package controller_test

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"testing"

	"github.com/golang/mock/gomock"

	"github.com/typical-go/typical-rest-server/internal/app/data_access/mysqldb"
	"github.com/typical-go/typical-rest-server/internal/app/domain/mymusic/controller"
	"github.com/typical-go/typical-rest-server/internal/app/domain/mymusic/service_mock"
	"github.com/typical-go/typical-rest-server/pkg/echotest"
	"github.com/typical-go/typical-rest-server/pkg/typrest"
)

type (
	SongCntrlFn       func(*service_mock.MockSongSvc)
	SongCntrlTestCase struct {
		TestName string
		echotest.TestCase
		SongCntrlFn
	}
)

func CreateSongCntrl(t *testing.T, fn SongCntrlFn) (*controller.SongCntrl, *gomock.Controller) {
	mock := gomock.NewController(t)
	mockSvc := service_mock.NewMockSongSvc(mock)
	if fn != nil {
		fn(mockSvc)
	}
	return &controller.SongCntrl{SongSvc: mockSvc}, mock
}

func TestBookController_RetrieveOne(t *testing.T) {
	testcases := []SongCntrlTestCase{
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
					Body: "{\"id\":1,\"title\":\"title1\",\"artist\":\"artist1\",\"update_at\":\"0001-01-01T00:00:00Z\",\"created_at\":\"0001-01-01T00:00:00Z\"}\n",
				},
			},
			SongCntrlFn: func(svc *service_mock.MockSongSvc) {
				svc.EXPECT().FindOne(gomock.Any(), "1").Return(&mysqldb.Song{ID: 1, Title: "title1", Artist: "artist1"}, nil)
			},
		},
		{
			TestName: "entity not found",
			TestCase: echotest.TestCase{
				Request: echotest.Request{
					Method:    http.MethodGet,
					Target:    "/",
					URLParams: map[string]string{"id": "3"},
				},
				ExpectedError: "code=404, message=Not Found",
			},
			SongCntrlFn: func(svc *service_mock.MockSongSvc) {
				svc.EXPECT().FindOne(gomock.Any(), "3").Return(nil, sql.ErrNoRows)
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
			SongCntrlFn: func(svc *service_mock.MockSongSvc) {
				svc.EXPECT().FindOne(gomock.Any(), "2").Return(nil, typrest.NewValidErr("some-validation"))
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
			SongCntrlFn: func(svc *service_mock.MockSongSvc) {
				svc.EXPECT().FindOne(gomock.Any(), "2").Return(nil, errors.New("some-error"))
			},
		},
	}

	for _, tt := range testcases {
		t.Run(tt.TestName, func(t *testing.T) {
			cntrl, mock := CreateSongCntrl(t, tt.SongCntrlFn)
			defer mock.Finish()
			tt.Execute(t, cntrl.FindOne)
		})
	}
}

func TestBookController_Retrieve(t *testing.T) {
	testcases := []SongCntrlTestCase{
		{
			TestCase: echotest.TestCase{
				Request: echotest.Request{
					Method: http.MethodGet,
					Target: "/",
				},
				ExpectedResponse: echotest.Response{
					Code: http.StatusOK,
					Body: "[{\"id\":1,\"title\":\"title1\",\"artist\":\"artist1\",\"update_at\":\"0001-01-01T00:00:00Z\",\"created_at\":\"0001-01-01T00:00:00Z\"},{\"id\":2,\"title\":\"title2\",\"artist\":\"artist2\",\"update_at\":\"0001-01-01T00:00:00Z\",\"created_at\":\"0001-01-01T00:00:00Z\"}]\n",
				},
			},
			SongCntrlFn: func(svc *service_mock.MockSongSvc) {
				svc.EXPECT().Find(gomock.Any()).Return([]*mysqldb.Song{
					&mysqldb.Song{ID: 1, Title: "title1", Artist: "artist1"},
					&mysqldb.Song{ID: 2, Title: "title2", Artist: "artist2"},
				}, nil)
			},
		},
		{
			TestCase: echotest.TestCase{
				Request: echotest.Request{
					Method: http.MethodGet,
					Target: "/",
				},
				ExpectedError: "code=500, message=some-error",
			},
			SongCntrlFn: func(svc *service_mock.MockSongSvc) {
				svc.EXPECT().Find(gomock.Any()).Return(nil, fmt.Errorf("some-error"))
			},
		},
	}

	for _, tt := range testcases {
		t.Run(tt.TestName, func(t *testing.T) {
			cntrl, mock := CreateSongCntrl(t, tt.SongCntrlFn)
			defer mock.Finish()
			tt.Execute(t, cntrl.Find)
		})
	}
}

func TestBookController_Update(t *testing.T) {
	testcases := []SongCntrlTestCase{
		{
			TestCase: echotest.TestCase{
				Request: echotest.Request{
					Method:    http.MethodPut,
					Target:    "/",
					URLParams: map[string]string{"id": "1"},
					Header:    echotest.HeaderForJSON(),
					Body:      `{"title":"some-title", "artist": "some-artist"}`,
				},
				ExpectedError: "code=500, message=some-error",
			},
			SongCntrlFn: func(svc *service_mock.MockSongSvc) {
				svc.EXPECT().
					Update(gomock.Any(), "1", &mysqldb.Song{ID: 1, Title: "some-title", Artist: "some-artist"}).
					Return(nil, errors.New("some-error"))
			},
		},
		{
			TestCase: echotest.TestCase{
				Request: echotest.Request{
					Method:    http.MethodPut,
					Target:    "/",
					URLParams: map[string]string{"id": "1"},
					Header:    echotest.HeaderForJSON(),
					Body:      `{"title":"some-title", "artist": "some-artist"}`,
				},
				ExpectedResponse: echotest.Response{
					Code: http.StatusOK,
					Body: "{\"id\":1,\"title\":\"some-title\",\"artist\":\"some-artist\",\"update_at\":\"0001-01-01T00:00:00Z\",\"created_at\":\"0001-01-01T00:00:00Z\"}\n",
				},
			},
			SongCntrlFn: func(svc *service_mock.MockSongSvc) {
				svc.EXPECT().
					Update(gomock.Any(), "1", &mysqldb.Song{ID: 1, Title: "some-title", Artist: "some-artist"}).
					Return(&mysqldb.Song{ID: 1, Title: "some-title", Artist: "some-artist"}, nil)
			},
		},
	}

	for _, tt := range testcases {
		t.Run(tt.TestName, func(t *testing.T) {
			cntrl, mock := CreateSongCntrl(t, tt.SongCntrlFn)
			defer mock.Finish()
			tt.Execute(t, cntrl.Update)
		})
	}
}

func TestBookController_Patch(t *testing.T) {
	testcases := []SongCntrlTestCase{
		{
			TestCase: echotest.TestCase{
				Request: echotest.Request{
					Method:    http.MethodPut,
					Target:    "/",
					URLParams: map[string]string{"id": "1"},
					Header:    echotest.HeaderForJSON(),
					Body:      `{"title":"some-title", "artist": "some-artist"}`,
				},
				ExpectedError: "code=500, message=some-error",
			},
			SongCntrlFn: func(svc *service_mock.MockSongSvc) {
				svc.EXPECT().
					Patch(gomock.Any(), "1", &mysqldb.Song{ID: 1, Title: "some-title", Artist: "some-artist"}).
					Return(nil, errors.New("some-error"))
			},
		},
		{
			TestCase: echotest.TestCase{
				Request: echotest.Request{
					Method:    http.MethodPut,
					Target:    "/",
					URLParams: map[string]string{"id": "1"},
					Header:    echotest.HeaderForJSON(),
					Body:      `{"title":"some-title", "artist": "some-artist"}`,
				},
				ExpectedResponse: echotest.Response{
					Code: http.StatusOK,
					Body: "{\"id\":1,\"title\":\"some-title\",\"artist\":\"some-artist\",\"update_at\":\"0001-01-01T00:00:00Z\",\"created_at\":\"0001-01-01T00:00:00Z\"}\n",
				},
			},
			SongCntrlFn: func(svc *service_mock.MockSongSvc) {
				svc.EXPECT().
					Patch(gomock.Any(), "1", &mysqldb.Song{ID: 1, Title: "some-title", Artist: "some-artist"}).
					Return(&mysqldb.Song{ID: 1, Title: "some-title", Artist: "some-artist"}, nil)
			},
		},
	}

	for _, tt := range testcases {
		t.Run(tt.TestName, func(t *testing.T) {
			cntrl, mock := CreateSongCntrl(t, tt.SongCntrlFn)
			defer mock.Finish()
			tt.Execute(t, cntrl.Patch)
		})
	}
}

func TestBookController_Delete(t *testing.T) {
	testcases := []SongCntrlTestCase{
		{
			TestCase: echotest.TestCase{
				Request: echotest.Request{
					Method:    http.MethodDelete,
					Target:    "/",
					URLParams: map[string]string{"id": "1"},
				},
				ExpectedResponse: echotest.Response{
					Code: http.StatusNoContent,
				},
			},
			SongCntrlFn: func(svc *service_mock.MockSongSvc) {
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
			SongCntrlFn: func(svc *service_mock.MockSongSvc) {
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
			SongCntrlFn: func(svc *service_mock.MockSongSvc) {
				svc.EXPECT().Delete(gomock.Any(), "1").Return(typrest.NewValidErr("some-validation"))
			},
		},
	}

	for _, tt := range testcases {
		t.Run(tt.TestName, func(t *testing.T) {
			cntrl, mock := CreateSongCntrl(t, tt.SongCntrlFn)
			defer mock.Finish()
			tt.Execute(t, cntrl.Delete)
		})
	}
}

func TestBookController_Create(t *testing.T) {
	testcases := []SongCntrlTestCase{
		{
			TestCase: echotest.TestCase{
				Request: echotest.Request{
					Method: http.MethodPost,
					Target: "/",
					Body:   `invalid}`,
					Header: echotest.HeaderForJSON(),
				},
				ExpectedError: "code=400, message=Syntax error: offset=1, error=invalid character 'i' looking for beginning of value, internal=invalid character 'i' looking for beginning of value",
			},
		},
		{
			TestCase: echotest.TestCase{
				Request: echotest.Request{
					Method: http.MethodPost,
					Target: "/",
					Body:   `{"artist":"some-artist", "title":"some-title"}`,
					Header: echotest.HeaderForJSON(),
				},
				ExpectedError: "code=500, message=some-error",
			},
			SongCntrlFn: func(svc *service_mock.MockSongSvc) {
				svc.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf("some-error"))
			},
		},
		{
			TestCase: echotest.TestCase{
				Request: echotest.Request{
					Method: http.MethodPost,
					Target: "/",
					Body:   `{"artist":"some-artist", "title":"some-title"}`,
					Header: echotest.HeaderForJSON(),
				},
				ExpectedResponse: echotest.Response{
					Body:   "{\"id\":999,\"title\":\"some-title\",\"artist\":\"some-artist\",\"update_at\":\"0001-01-01T00:00:00Z\",\"created_at\":\"0001-01-01T00:00:00Z\"}\n",
					Code:   http.StatusCreated,
					Header: map[string]string{"Location": "/songs/999"},
				},
			},
			SongCntrlFn: func(svc *service_mock.MockSongSvc) {
				svc.EXPECT().
					Create(gomock.Any(), gomock.Any()).
					Return(&mysqldb.Song{ID: 999, Artist: "some-artist", Title: "some-title"}, nil)
			},
		},
	}

	for _, tt := range testcases {
		t.Run(tt.TestName, func(t *testing.T) {
			cntrl, mock := CreateSongCntrl(t, tt.SongCntrlFn)
			defer mock.Finish()
			tt.Execute(t, cntrl.Create)
		})
	}
}
