package mysqldb_test

import (
	"context"
	"database/sql"
	"errors"
	"regexp"
	"testing"
	"time"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"
	"github.com/typical-go/typical-rest-server/internal/app/data_access/mysqldb"
	"github.com/typical-go/typical-rest-server/pkg/dbkit"
	"github.com/typical-go/typical-rest-server/pkg/dbtxn"
)

type bookRepoFn func(sqlmock.Sqlmock)

func createSongRepo(fn bookRepoFn) (mysqldb.SongRepo, *sql.DB) {
	db, mock, _ := sqlmock.New()
	if fn != nil {
		fn(mock)
	}
	return &mysqldb.SongRepoImpl{DB: db}, db
}

func TestSongRepoImpl_Create(t *testing.T) {
	testcases := []struct {
		testName           string
		book               *mysqldb.Song
		bookRepoFn         bookRepoFn
		expectedInsertedID int64
		expectedErr        string
	}{
		{
			testName:    "begin error",
			book:        &mysqldb.Song{Title: "some-title", Artist: "some-artist"},
			expectedErr: "dbtxn: some-error",
			bookRepoFn: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin().WillReturnError(errors.New("some-error"))
			},
		},
		{
			testName:    "insert error",
			book:        &mysqldb.Song{Title: "some-title", Artist: "some-artist"},
			expectedErr: "some-error",
			bookRepoFn: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO songs (title,artist,created_at,updated_at) VALUES (?,?,?,?) RETURNING "id"`)).
					WithArgs("some-title", "some-artist", sqlmock.AnyArg(), sqlmock.AnyArg()).
					WillReturnError(errors.New("some-error"))
			},
		},
		{
			book: &mysqldb.Song{Title: "some-title", Artist: "some-artist"},
			bookRepoFn: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO songs (title,artist,created_at,updated_at) VALUES (?,?,?,?) RETURNING "id"`)).
					WithArgs("some-title", "some-artist", sqlmock.AnyArg(), sqlmock.AnyArg()).
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(999))
			},
			expectedInsertedID: 999,
		},
	}

	for _, tt := range testcases {
		t.Run(tt.testName, func(t *testing.T) {
			repo, db := createSongRepo(tt.bookRepoFn)
			defer db.Close()

			ctx := context.Background()
			dbtxn.Begin(&ctx)

			id, err := repo.Create(ctx, tt.book)

			if tt.expectedErr != "" {
				require.EqualError(t, err, tt.expectedErr)
				require.EqualError(t, dbtxn.Error(ctx), tt.expectedErr)
			} else {
				require.NoError(t, err)
				require.NoError(t, dbtxn.Error(ctx))
				require.Equal(t, tt.expectedInsertedID, id)
			}

		})
	}
}

func TestSongRepoImpl_Update(t *testing.T) {
	testcases := []struct {
		testName            string
		book                *mysqldb.Song
		bookRepoFn          bookRepoFn
		opt                 dbkit.UpdateOption
		expectedErr         string
		expectedAffectedRow int64
	}{
		{
			testName:            "update error",
			book:                &mysqldb.Song{Title: "new-title", Artist: "new-artist"},
			opt:                 dbkit.Equal(mysqldb.SongTable.ID, 888),
			expectedErr:         "dbtxn: begin-error",
			expectedAffectedRow: -1,
			bookRepoFn: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin().WillReturnError(errors.New("begin-error"))
			},
		},
		{
			testName: "update error",
			book:     &mysqldb.Song{Title: "new-title", Artist: "new-artist"},
			opt:      dbkit.Equal(mysqldb.SongTable.ID, 888),
			bookRepoFn: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta(`UPDATE songs SET title = ?, artist = ?, updated_at = ? WHERE id = ?`)).
					WithArgs("new-title", "new-artist", sqlmock.AnyArg(), 888).
					WillReturnError(errors.New("some-update-error"))
			},
			expectedErr:         "some-update-error",
			expectedAffectedRow: -1,
		},
		{
			testName: "complete book",
			book:     &mysqldb.Song{Title: "new-title", Artist: "new-artist"},
			opt:      dbkit.Equal(mysqldb.SongTable.ID, 888),
			bookRepoFn: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta(`UPDATE songs SET title = ?, artist = ?, updated_at = ? WHERE id = ?`)).
					WithArgs("new-title", "new-artist", sqlmock.AnyArg(), 888).
					WillReturnResult(sqlmock.NewResult(0, 1))
			},
			expectedAffectedRow: 1,
		},
		{
			testName: "empty artist",
			book:     &mysqldb.Song{Title: "new-title"},
			opt:      dbkit.Equal(mysqldb.SongTable.ID, 888),
			bookRepoFn: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta(`UPDATE songs SET title = ?, artist = ?, updated_at = ? WHERE id = ?`)).
					WithArgs("new-title", "", sqlmock.AnyArg(), 888).
					WillReturnResult(sqlmock.NewResult(0, 1))
			},
			expectedAffectedRow: 1,
		},
		{
			testName: "empty title",
			book:     &mysqldb.Song{Artist: "new-artist"},
			opt:      dbkit.Equal(mysqldb.SongTable.ID, 888),
			bookRepoFn: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta(`UPDATE songs SET title = ?, artist = ?, updated_at = ? WHERE id = ?`)).
					WithArgs("", "new-artist", sqlmock.AnyArg(), 888).
					WillReturnResult(sqlmock.NewResult(0, 1))
			},
			expectedAffectedRow: 1,
		},
	}
	for _, tt := range testcases {
		t.Run(tt.testName, func(t *testing.T) {
			repo, db := createSongRepo(tt.bookRepoFn)
			defer db.Close()

			ctx := context.Background()
			dbtxn.Begin(&ctx)

			affectedRow, err := repo.Update(ctx, tt.book, tt.opt)
			if tt.expectedErr != "" {
				require.EqualError(t, err, tt.expectedErr)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.expectedAffectedRow, affectedRow)
			}
		})
	}
}

func TestSongRepoImpl_Patch(t *testing.T) {
	testcases := []struct {
		testName            string
		book                *mysqldb.Song
		bookRepoFn          bookRepoFn
		opt                 dbkit.UpdateOption
		expectedErr         string
		expectedAffectedRow int64
	}{
		{
			testName: "begin error",
			book:     &mysqldb.Song{Title: "new-title", Artist: "new-artist"},
			opt:      dbkit.Equal(mysqldb.SongTable.ID, 888),
			bookRepoFn: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin().WillReturnError(errors.New("begin-error"))
			},
			expectedErr:         "dbtxn: begin-error",
			expectedAffectedRow: -1,
		},
		{
			testName: "update error",
			book:     &mysqldb.Song{Title: "new-title", Artist: "new-artist"},
			opt:      dbkit.Equal(mysqldb.SongTable.ID, 888),
			bookRepoFn: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta(`UPDATE songs SET title = ?, artist = ?, updated_at = ? WHERE id = ?`)).
					WithArgs("new-title", "new-artist", sqlmock.AnyArg(), 888).
					WillReturnError(errors.New("some-update-error"))
			},
			expectedErr:         "some-update-error",
			expectedAffectedRow: -1,
		},
		{
			testName: "complete book",
			book:     &mysqldb.Song{Title: "new-title", Artist: "new-artist"},
			opt:      dbkit.Equal(mysqldb.SongTable.ID, 888),
			bookRepoFn: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta(`UPDATE songs SET title = ?, artist = ?, updated_at = ? WHERE id = ?`)).
					WithArgs("new-title", "new-artist", sqlmock.AnyArg(), 888).
					WillReturnResult(sqlmock.NewResult(0, 1))
			},
			expectedAffectedRow: 1,
		},
		{
			testName: "empty artist",
			book:     &mysqldb.Song{Title: "new-title"},
			opt:      dbkit.Equal(mysqldb.SongTable.ID, 888),
			bookRepoFn: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta(`UPDATE songs SET title = ?, updated_at = ? WHERE id = ?`)).
					WithArgs("new-title", sqlmock.AnyArg(), 888).
					WillReturnResult(sqlmock.NewResult(0, 1))
			},
			expectedAffectedRow: 1,
		},
		{
			testName: "empty title",
			book:     &mysqldb.Song{Artist: "new-artist"},
			opt:      dbkit.Equal(mysqldb.SongTable.ID, 888),
			bookRepoFn: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta(`UPDATE songs SET artist = ?, updated_at = ? WHERE id = ?`)).
					WithArgs("new-artist", sqlmock.AnyArg(), 888).
					WillReturnResult(sqlmock.NewResult(0, 1))
			},
			expectedAffectedRow: 1,
		},
	}
	for _, tt := range testcases {
		t.Run(tt.testName, func(t *testing.T) {
			repo, db := createSongRepo(tt.bookRepoFn)
			defer db.Close()

			ctx := context.Background()
			dbtxn.Begin(&ctx)

			affectedRow, err := repo.Patch(ctx, tt.book, tt.opt)
			if tt.expectedErr != "" {
				require.EqualError(t, err, tt.expectedErr)
				require.EqualError(t, dbtxn.Error(ctx), tt.expectedErr)
			} else {
				require.NoError(t, err)
				require.NoError(t, dbtxn.Error(ctx))
				require.Equal(t, tt.expectedAffectedRow, affectedRow)
			}
		})
	}
}

func TestSongRepoImpl_Retrieve(t *testing.T) {
	now := time.Now()
	testcases := []struct {
		testName    string
		opts        []dbkit.SelectOption
		expected    []*mysqldb.Song
		expectedErr string
		bookRepoFn  bookRepoFn
	}{
		{

			opts:        []dbkit.SelectOption{},
			expected:    []*mysqldb.Song{},
			expectedErr: "some-error",
			bookRepoFn: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT id, title, artist, updated_at, created_at FROM songs`).
					WillReturnError(errors.New("some-error"))
			},
		},
		{

			opts: []dbkit.SelectOption{},
			expected: []*mysqldb.Song{
				&mysqldb.Song{ID: 1234, Title: "some-title4", Artist: "some-artist4", UpdatedAt: now, CreatedAt: now},
				&mysqldb.Song{ID: 1235, Title: "some-title5", Artist: "some-artist5", UpdatedAt: now, CreatedAt: now},
			},
			bookRepoFn: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT id, title, artist, updated_at, created_at FROM songs`).
					WillReturnRows(sqlmock.
						NewRows([]string{"id", "title", "artist", "updated_at", "created_at"}).
						AddRow("1234", "some-title4", "some-artist4", now, now).
						AddRow("1235", "some-title5", "some-artist5", now, now),
					)
			},
		},
	}
	for _, tt := range testcases {
		t.Run(tt.testName, func(t *testing.T) {
			repo, db := createSongRepo(tt.bookRepoFn)
			defer db.Close()

			songs, err := repo.Find(context.Background(), tt.opts...)
			if tt.expectedErr != "" {
				require.EqualError(t, err, tt.expectedErr)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.expected, songs)
			}
		})
	}
}

func TestSongRepoImpl_Delete(t *testing.T) {
	testcases := []struct {
		testName            string
		opt                 dbkit.DeleteOption
		bookRepoFn          bookRepoFn
		expectedErr         string
		expectedAffectedRow int64
	}{
		{
			testName:    "begin error",
			opt:         dbkit.Equal("id", 666),
			expectedErr: "dbtxn: begin-error",
			bookRepoFn: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin().WillReturnError(errors.New("begin-error"))
			},
		},
		{
			testName:    "delete error",
			opt:         dbkit.Equal("id", 666),
			expectedErr: "delete-error",
			bookRepoFn: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM songs WHERE id = ?`)).
					WithArgs(666).
					WillReturnError(errors.New("delete-error"))
			},
		},
		{
			opt: dbkit.Equal("id", 555),
			bookRepoFn: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM songs WHERE id = ?`)).
					WithArgs(555).
					WillReturnResult(sqlmock.NewResult(0, 1))
			},
			expectedAffectedRow: 1,
		},
	}
	for _, tt := range testcases {
		t.Run(tt.testName, func(t *testing.T) {
			repo, db := createSongRepo(tt.bookRepoFn)
			defer db.Close()

			ctx := context.Background()
			dbtxn.Begin(&ctx)

			affectedRow, err := repo.Delete(ctx, tt.opt)
			if tt.expectedErr != "" {
				require.EqualError(t, err, tt.expectedErr)
				require.EqualError(t, dbtxn.Error(ctx), tt.expectedErr)
			} else {
				require.NoError(t, err)
				require.NoError(t, dbtxn.Error(ctx))
				require.Equal(t, tt.expectedAffectedRow, affectedRow)
			}
		})
	}
}
