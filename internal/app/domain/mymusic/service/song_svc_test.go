package service_test

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"github.com/typical-go/typical-rest-server/internal/app/data_access/mysqldb"
	"github.com/typical-go/typical-rest-server/internal/app/data_access/mysqldb_mock"
	"github.com/typical-go/typical-rest-server/internal/app/domain/mymusic/service"
	"github.com/typical-go/typical-rest-server/pkg/dbkit"
)

type songSvcFn func(mockRepo *mysqldb_mock.MockSongRepo)

func createSongSvc(t *testing.T, fn songSvcFn) (*service.SongSvcImpl, *gomock.Controller) {
	mock := gomock.NewController(t)
	mockRepo := mysqldb_mock.NewMockSongRepo(mock)
	if fn != nil {
		fn(mockRepo)
	}

	return &service.SongSvcImpl{
		SongRepo: mockRepo,
	}, mock
}

func TestSongSvc_Create(t *testing.T) {
	testcases := []struct {
		testName    string
		songSvcFn   songSvcFn
		song        *mysqldb.Song
		expected    *mysqldb.Song
		expectedErr string
	}{
		{
			testName:    "validation error",
			song:        &mysqldb.Song{},
			expectedErr: "Key: 'Song.Title' Error:Field validation for 'Title' failed on the 'required' tag\nKey: 'Song.Artist' Error:Field validation for 'Artist' failed on the 'required' tag",
		},
		{
			testName:    "create error",
			song:        &mysqldb.Song{Artist: "some-artist", Title: "some-title"},
			expectedErr: "create-error",
			songSvcFn: func(mockRepo *mysqldb_mock.MockSongRepo) {
				mockRepo.EXPECT().
					Create(gomock.Any(), &mysqldb.Song{Artist: "some-artist", Title: "some-title"}).
					Return(int64(-1), errors.New("create-error"))
			},
		},
		{
			testName:    "Find error",
			song:        &mysqldb.Song{Artist: "some-artist", Title: "some-title"},
			expectedErr: "Find-error",
			songSvcFn: func(mockRepo *mysqldb_mock.MockSongRepo) {
				mockRepo.EXPECT().
					Create(gomock.Any(), &mysqldb.Song{Artist: "some-artist", Title: "some-title"}).
					Return(int64(1), nil)
				mockRepo.EXPECT().
					Find(gomock.Any(), dbkit.Equal("id", int64(1))).
					Return(nil, errors.New("Find-error"))
			},
		},
		{
			song: &mysqldb.Song{
				Artist: "some-artist",
				Title:  "some-title",
			},
			expected: &mysqldb.Song{Artist: "some-artist", Title: "some-title"},
			songSvcFn: func(mockRepo *mysqldb_mock.MockSongRepo) {
				mockRepo.EXPECT().
					Create(gomock.Any(), &mysqldb.Song{Artist: "some-artist", Title: "some-title"}).
					Return(int64(1), nil)
				mockRepo.EXPECT().
					Find(gomock.Any(), dbkit.Equal("id", int64(1))).
					Return([]*mysqldb.Song{{Artist: "some-artist", Title: "some-title"}}, nil)
			},
		},
	}

	for _, tt := range testcases {
		t.Run(tt.testName, func(t *testing.T) {
			svc, mock := createSongSvc(t, tt.songSvcFn)
			defer mock.Finish()

			id, err := svc.Create(context.Background(), tt.song)
			if tt.expectedErr != "" {
				require.EqualError(t, err, tt.expectedErr)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.expected, id)
			}
		})
	}
}

func TestSongSvc_RetrieveOne(t *testing.T) {
	testcases := []struct {
		testName    string
		songSvcFn   songSvcFn
		paramID     string
		expected    *mysqldb.Song
		expectedErr string
	}{
		{
			paramID:     "",
			expectedErr: "paramID is missing",
		},
		{
			paramID: "1",
			songSvcFn: func(mockRepo *mysqldb_mock.MockSongRepo) {
				mockRepo.EXPECT().
					Find(gomock.Any(), dbkit.Equal("id", int64(1))).
					Return(nil, errors.New("some-error"))
			},
			expectedErr: "some-error",
		},
		{
			paramID: "1",
			songSvcFn: func(mockRepo *mysqldb_mock.MockSongRepo) {
				mockRepo.EXPECT().
					Find(gomock.Any(), dbkit.Equal("id", int64(1))).
					Return([]*mysqldb.Song{
						{
							ID:    1,
							Title: "some-title",
						},
					}, nil)
			},
			expected: &mysqldb.Song{
				ID:    1,
				Title: "some-title",
			},
		},
	}

	for _, tt := range testcases {
		t.Run(tt.testName, func(t *testing.T) {
			svc, mock := createSongSvc(t, tt.songSvcFn)
			defer mock.Finish()

			song, err := svc.FindOne(context.Background(), tt.paramID)
			if tt.expectedErr != "" {
				require.EqualError(t, err, tt.expectedErr)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.expected, song)
			}
		})
	}
}

func TestSongSvc_Retrieve(t *testing.T) {
	testcases := []struct {
		testName    string
		songSvcFn   songSvcFn
		expected    []*mysqldb.Song
		expectedErr string
	}{}
	for _, tt := range testcases {
		t.Run(tt.testName, func(t *testing.T) {
			svc, mock := createSongSvc(t, tt.songSvcFn)
			defer mock.Finish()

			songs, err := svc.Find(context.Background())
			if tt.expectedErr != "" {
				require.EqualError(t, err, tt.expectedErr)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.expected, songs)
			}
		})
	}
}

func TestSongSvc_Delete(t *testing.T) {
	testcases := []struct {
		testName    string
		songSvcFn   songSvcFn
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
			songSvcFn: func(mockRepo *mysqldb_mock.MockSongRepo) {
				mockRepo.EXPECT().
					Delete(gomock.Any(), dbkit.Equal(mysqldb.SongTable.ID, int64(1))).
					Return(int64(0), errors.New("some-error"))
			},
		},
		{
			paramID: "1",
			songSvcFn: func(mockRepo *mysqldb_mock.MockSongRepo) {
				mockRepo.EXPECT().
					Delete(gomock.Any(), dbkit.Equal(mysqldb.SongTable.ID, int64(1))).
					Return(int64(1), nil)
			},
		},
		{
			testName: "success even if no affected row (idempotent)",
			paramID:  "1",
			songSvcFn: func(mockRepo *mysqldb_mock.MockSongRepo) {
				mockRepo.EXPECT().
					Delete(gomock.Any(), dbkit.Equal(mysqldb.SongTable.ID, int64(1))).
					Return(int64(0), nil)
			},
		},
	}
	for _, tt := range testcases {
		t.Run(tt.testName, func(t *testing.T) {
			svc, mock := createSongSvc(t, tt.songSvcFn)
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

func TestSongSvc_Update(t *testing.T) {
	testcases := []struct {
		testName    string
		songSvcFn   songSvcFn
		paramID     string
		song        *mysqldb.Song
		expected    *mysqldb.Song
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
			song:        &mysqldb.Song{},
			expectedErr: "Key: 'Song.Title' Error:Field validation for 'Title' failed on the 'required' tag\nKey: 'Song.Artist' Error:Field validation for 'Artist' failed on the 'required' tag",
		},
		{
			testName:    "update error",
			paramID:     "1",
			song:        &mysqldb.Song{Artist: "some-artist", Title: "some-title"},
			expectedErr: "update error",
			songSvcFn: func(mockRepo *mysqldb_mock.MockSongRepo) {
				mockRepo.EXPECT().
					Update(gomock.Any(), &mysqldb.Song{Artist: "some-artist", Title: "some-title"}, dbkit.Equal("id", int64(1))).
					Return(int64(-1), errors.New("update error"))
			},
		},
		{
			testName:    "nothing to update",
			paramID:     "1",
			song:        &mysqldb.Song{Artist: "some-artist", Title: "some-title"},
			expectedErr: "sql: no rows in result set",
			songSvcFn: func(mockRepo *mysqldb_mock.MockSongRepo) {
				mockRepo.EXPECT().
					Update(gomock.Any(), &mysqldb.Song{Artist: "some-artist", Title: "some-title"}, dbkit.Equal("id", int64(1))).
					Return(int64(0), nil)
			},
		},
		{
			testName:    "Find error",
			paramID:     "1",
			song:        &mysqldb.Song{Artist: "some-artist", Title: "some-title"},
			expectedErr: "Find-error",
			songSvcFn: func(mockRepo *mysqldb_mock.MockSongRepo) {
				mockRepo.EXPECT().
					Update(gomock.Any(), &mysqldb.Song{Artist: "some-artist", Title: "some-title"}, dbkit.Equal("id", int64(1))).
					Return(int64(1), nil)
				mockRepo.EXPECT().
					Find(gomock.Any(), dbkit.Equal("id", int64(1))).
					Return(nil, errors.New("Find-error"))
			},
		},
		{
			testName:    "Find error",
			paramID:     "1",
			song:        &mysqldb.Song{Artist: "some-artist", Title: "some-title"},
			expectedErr: "Find-error",
			songSvcFn: func(mockRepo *mysqldb_mock.MockSongRepo) {
				mockRepo.EXPECT().
					Update(gomock.Any(), &mysqldb.Song{Artist: "some-artist", Title: "some-title"}, dbkit.Equal("id", int64(1))).
					Return(int64(1), nil)
				mockRepo.EXPECT().
					Find(gomock.Any(), dbkit.Equal("id", int64(1))).
					Return(nil, errors.New("Find-error"))
			},
		},
	}
	for _, tt := range testcases {
		t.Run(tt.testName, func(t *testing.T) {
			svc, mock := createSongSvc(t, tt.songSvcFn)
			defer mock.Finish()

			song, err := svc.Update(context.Background(), tt.paramID, tt.song)
			if tt.expectedErr != "" {
				require.EqualError(t, err, tt.expectedErr)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.expected, song)
			}
		})
	}
}

func TestSongSvc_Patch(t *testing.T) {
	testcases := []struct {
		testName    string
		songSvcFn   songSvcFn
		paramID     string
		song        *mysqldb.Song
		expected    *mysqldb.Song
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
			song:        &mysqldb.Song{Artist: "some-artist", Title: "some-title"},
			expectedErr: "patch-error",
			songSvcFn: func(mockRepo *mysqldb_mock.MockSongRepo) {
				mockRepo.EXPECT().
					Patch(gomock.Any(), &mysqldb.Song{Artist: "some-artist", Title: "some-title"}, dbkit.Equal("id", int64(1))).
					Return(int64(-1), errors.New("patch-error"))
			},
		},
		{
			testName:    "patch error",
			paramID:     "1",
			song:        &mysqldb.Song{Artist: "some-artist", Title: "some-title"},
			expectedErr: "sql: no rows in result set",
			songSvcFn: func(mockRepo *mysqldb_mock.MockSongRepo) {
				mockRepo.EXPECT().
					Patch(gomock.Any(), &mysqldb.Song{Artist: "some-artist", Title: "some-title"}, dbkit.Equal("id", int64(1))).
					Return(int64(0), nil)
			},
		},
		{
			testName:    "Find error",
			paramID:     "1",
			song:        &mysqldb.Song{Artist: "some-artist", Title: "some-title"},
			expectedErr: "Find-error",
			songSvcFn: func(mockRepo *mysqldb_mock.MockSongRepo) {
				mockRepo.EXPECT().
					Patch(gomock.Any(), &mysqldb.Song{Artist: "some-artist", Title: "some-title"}, dbkit.Equal("id", int64(1))).
					Return(int64(1), nil)
				mockRepo.EXPECT().
					Find(gomock.Any(), dbkit.Equal("id", int64(1))).
					Return(nil, errors.New("Find-error"))
			},
		},
		{
			paramID:  "1",
			song:     &mysqldb.Song{Artist: "some-artist", Title: "some-title"},
			expected: &mysqldb.Song{Artist: "some-artist", Title: "some-title"},
			songSvcFn: func(mockRepo *mysqldb_mock.MockSongRepo) {
				mockRepo.EXPECT().
					Patch(gomock.Any(), &mysqldb.Song{Artist: "some-artist", Title: "some-title"}, dbkit.Equal("id", int64(1))).
					Return(int64(1), nil)
				mockRepo.EXPECT().
					Find(gomock.Any(), dbkit.Equal("id", int64(1))).
					Return([]*mysqldb.Song{{Artist: "some-artist", Title: "some-title"}}, nil)
			},
		},
	}
	for _, tt := range testcases {
		t.Run(tt.testName, func(t *testing.T) {
			svc, mock := createSongSvc(t, tt.songSvcFn)
			defer mock.Finish()

			song, err := svc.Patch(context.Background(), tt.paramID, tt.song)
			if tt.expectedErr != "" {
				require.EqualError(t, err, tt.expectedErr)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.expected, song)
			}
		})
	}
}
