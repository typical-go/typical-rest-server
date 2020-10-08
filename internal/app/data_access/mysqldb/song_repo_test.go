package mysqldb_test

import (
	"context"
	"database/sql"
	"errors"
	"regexp"
	"testing"
	"time"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	sq "github.com/Masterminds/squirrel"
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
	return mysqldb.NewSongRepo(mysqldb.SongRepoImpl{DB: db}), db
}

func TestSongRepoImpl_Create(t *testing.T) {
	testcases := []struct {
		TestName    string
		Song        *mysqldb.Song
		SongRepoFn  bookRepoFn
		Expected    int64
		ExpectedErr string
	}{
		{
			TestName:    "begin error",
			Song:        &mysqldb.Song{Title: "some-title", Artist: "some-artist"},
			ExpectedErr: "dbtxn: some-error",
			SongRepoFn: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin().WillReturnError(errors.New("some-error"))
			},
		},
		{
			TestName:    "insert error",
			Song:        &mysqldb.Song{Title: "some-title", Artist: "some-artist"},
			ExpectedErr: "some-error",
			SongRepoFn: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO songs (title,artist,created_at,updated_at) VALUES (?,?,?,?)`)).
					WithArgs("some-title", "some-artist", sqlmock.AnyArg(), sqlmock.AnyArg()).
					WillReturnError(errors.New("some-error"))
			},
		},
		{
			Song: &mysqldb.Song{Title: "some-title", Artist: "some-artist"},
			SongRepoFn: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO songs (title,artist,created_at,updated_at) VALUES (?,?,?,?)`)).
					WithArgs("some-title", "some-artist", sqlmock.AnyArg(), sqlmock.AnyArg()).
					WillReturnResult(sqlmock.NewResult(999, 1))
			},
			Expected: 999,
		},
	}

	for _, tt := range testcases {
		t.Run(tt.TestName, func(t *testing.T) {
			repo, db := createSongRepo(tt.SongRepoFn)
			defer db.Close()

			ctx := context.Background()
			dbtxn.Begin(&ctx)

			id, err := repo.Create(ctx, tt.Song)

			if tt.ExpectedErr != "" {
				require.EqualError(t, err, tt.ExpectedErr)
				require.EqualError(t, dbtxn.Error(ctx), tt.ExpectedErr)
			} else {
				require.NoError(t, err)
				require.NoError(t, dbtxn.Error(ctx))
				require.Equal(t, tt.Expected, id)
			}

		})
	}
}

func TestSongRepoImpl_Update(t *testing.T) {
	testcases := []struct {
		TestName    string
		Song        *mysqldb.Song
		SongRepoFn  bookRepoFn
		Opt         dbkit.UpdateOption
		ExpectedErr string
		Expected    int64
	}{
		{
			TestName:    "update error",
			Song:        &mysqldb.Song{Title: "new-title", Artist: "new-artist"},
			Opt:         dbkit.Equal(mysqldb.SongTable.ID, 888),
			ExpectedErr: "dbtxn: begin-error",
			Expected:    -1,
			SongRepoFn: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin().WillReturnError(errors.New("begin-error"))
			},
		},
		{
			TestName: "update error",
			Song:     &mysqldb.Song{Title: "new-title", Artist: "new-artist"},
			Opt:      dbkit.Equal(mysqldb.SongTable.ID, 888),
			SongRepoFn: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta(`UPDATE songs SET title = ?, artist = ?, updated_at = ? WHERE id = ?`)).
					WithArgs("new-title", "new-artist", sqlmock.AnyArg(), 888).
					WillReturnError(errors.New("some-update-error"))
			},
			ExpectedErr: "some-update-error",
			Expected:    -1,
		},
		{
			TestName: "bad update option",
			Song:     &mysqldb.Song{Title: "new-title", Artist: "new-artist"},
			Opt: dbkit.NewUpdateOption(func(b sq.UpdateBuilder) (sq.UpdateBuilder, error) {
				return b, errors.New("bad-option")
			}),
			SongRepoFn: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta(`UPDATE songs SET title = ?, artist = ?, updated_at = ? WHERE id = ?`)).
					WithArgs("new-title", "new-artist", sqlmock.AnyArg(), 888).
					WillReturnResult(sqlmock.NewResult(0, 1))
			},
			ExpectedErr: "bad-option",
		},
		{
			TestName: "success",
			Song:     &mysqldb.Song{Title: "new-title", Artist: "new-artist"},
			Opt:      dbkit.Equal(mysqldb.SongTable.ID, 888),
			SongRepoFn: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta(`UPDATE songs SET title = ?, artist = ?, updated_at = ? WHERE id = ?`)).
					WithArgs("new-title", "new-artist", sqlmock.AnyArg(), 888).
					WillReturnResult(sqlmock.NewResult(0, 1))
			},
			Expected: 1,
		},
		{
			TestName: "success empty artist",
			Song:     &mysqldb.Song{Title: "new-title"},
			Opt:      dbkit.Equal(mysqldb.SongTable.ID, 888),
			SongRepoFn: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta(`UPDATE songs SET title = ?, artist = ?, updated_at = ? WHERE id = ?`)).
					WithArgs("new-title", "", sqlmock.AnyArg(), 888).
					WillReturnResult(sqlmock.NewResult(0, 1))
			},
			Expected: 1,
		},
		{
			TestName: "success empty title",
			Song:     &mysqldb.Song{Artist: "new-artist"},
			Opt:      dbkit.Equal(mysqldb.SongTable.ID, 888),
			SongRepoFn: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta(`UPDATE songs SET title = ?, artist = ?, updated_at = ? WHERE id = ?`)).
					WithArgs("", "new-artist", sqlmock.AnyArg(), 888).
					WillReturnResult(sqlmock.NewResult(0, 1))
			},
			Expected: 1,
		},
	}
	for _, tt := range testcases {
		t.Run(tt.TestName, func(t *testing.T) {
			repo, db := createSongRepo(tt.SongRepoFn)
			defer db.Close()

			ctx := context.Background()
			dbtxn.Begin(&ctx)

			affectedRow, err := repo.Update(ctx, tt.Song, tt.Opt)
			if tt.ExpectedErr != "" {
				require.EqualError(t, err, tt.ExpectedErr)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.Expected, affectedRow)
			}
		})
	}
}

func TestSongRepoImpl_Patch(t *testing.T) {
	testcases := []struct {
		TestName    string
		Song        *mysqldb.Song
		SongRepoFn  bookRepoFn
		Opt         dbkit.UpdateOption
		ExpectedErr string
		Expected    int64
	}{
		{
			TestName: "begin error",
			Song:     &mysqldb.Song{Title: "new-title", Artist: "new-artist"},
			Opt:      dbkit.Equal(mysqldb.SongTable.ID, 888),
			SongRepoFn: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin().WillReturnError(errors.New("begin-error"))
			},
			ExpectedErr: "dbtxn: begin-error",
			Expected:    -1,
		},
		{
			TestName: "update error",
			Song:     &mysqldb.Song{Title: "new-title", Artist: "new-artist"},
			Opt:      dbkit.Equal(mysqldb.SongTable.ID, 888),
			SongRepoFn: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta(`UPDATE songs SET title = ?, artist = ?, updated_at = ? WHERE id = ?`)).
					WithArgs("new-title", "new-artist", sqlmock.AnyArg(), 888).
					WillReturnError(errors.New("some-update-error"))
			},
			ExpectedErr: "some-update-error",
			Expected:    -1,
		},

		{
			TestName: "bad update option",
			Song:     &mysqldb.Song{Title: "new-title", Artist: "new-artist"},
			Opt: dbkit.NewUpdateOption(func(b sq.UpdateBuilder) (sq.UpdateBuilder, error) {
				return b, errors.New("bad-option")
			}),
			SongRepoFn: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta(`UPDATE songs SET title = ?, artist = ?, updated_at = ? WHERE id = ?`)).
					WithArgs("new-title", "new-artist", sqlmock.AnyArg(), 888).
					WillReturnResult(sqlmock.NewResult(0, 1))
			},
			ExpectedErr: "bad-option",
		},
		{
			TestName: "success",
			Song:     &mysqldb.Song{Title: "new-title", Artist: "new-artist"},
			Opt:      dbkit.Equal(mysqldb.SongTable.ID, 888),
			SongRepoFn: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta(`UPDATE songs SET title = ?, artist = ?, updated_at = ? WHERE id = ?`)).
					WithArgs("new-title", "new-artist", sqlmock.AnyArg(), 888).
					WillReturnResult(sqlmock.NewResult(0, 1))
			},
			Expected: 1,
		},
		{
			TestName: "success empty artist",
			Song:     &mysqldb.Song{Title: "new-title"},
			Opt:      dbkit.Equal(mysqldb.SongTable.ID, 888),
			SongRepoFn: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta(`UPDATE songs SET title = ?, updated_at = ? WHERE id = ?`)).
					WithArgs("new-title", sqlmock.AnyArg(), 888).
					WillReturnResult(sqlmock.NewResult(0, 1))
			},
			Expected: 1,
		},
		{
			TestName: "success empty title",
			Song:     &mysqldb.Song{Artist: "new-artist"},
			Opt:      dbkit.Equal(mysqldb.SongTable.ID, 888),
			SongRepoFn: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta(`UPDATE songs SET artist = ?, updated_at = ? WHERE id = ?`)).
					WithArgs("new-artist", sqlmock.AnyArg(), 888).
					WillReturnResult(sqlmock.NewResult(0, 1))
			},
			Expected: 1,
		},
	}
	for _, tt := range testcases {
		t.Run(tt.TestName, func(t *testing.T) {
			repo, db := createSongRepo(tt.SongRepoFn)
			defer db.Close()

			ctx := context.Background()
			dbtxn.Begin(&ctx)

			affectedRow, err := repo.Patch(ctx, tt.Song, tt.Opt)
			if tt.ExpectedErr != "" {
				require.EqualError(t, err, tt.ExpectedErr)
				require.EqualError(t, dbtxn.Error(ctx), tt.ExpectedErr)
			} else {
				require.NoError(t, err)
				require.NoError(t, dbtxn.Error(ctx))
				require.Equal(t, tt.Expected, affectedRow)
			}
		})
	}
}

func TestSongRepoImpl_Retrieve(t *testing.T) {
	now := time.Now()
	testcases := []struct {
		TestName    string
		Opts        []dbkit.SelectOption
		Expected    []*mysqldb.Song
		ExpectedErr string
		SongRepoFn  bookRepoFn
	}{
		{
			TestName: "sql error",
			Opts:     []dbkit.SelectOption{},
			Expected: []*mysqldb.Song{},
			SongRepoFn: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT id, title, artist, updated_at, created_at FROM songs`).
					WillReturnError(errors.New("some-error"))
			},
			ExpectedErr: "some-error",
		},
		{
			TestName: "bad option",
			Opts: []dbkit.SelectOption{
				dbkit.NewSelectOption(func(b sq.SelectBuilder) (sq.SelectBuilder, error) {
					return b, errors.New("bad-option")
				}),
			},
			SongRepoFn: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT id, title, artist, updated_at, created_at FROM songs`).
					WillReturnRows(sqlmock.
						NewRows([]string{"id", "title", "artist", "updated_at", "created_at"}).
						AddRow("1234", "some-title4", "some-artist4", now, now).
						AddRow("1235", "some-title5", "some-artist5", now, now),
					)
			},
			ExpectedErr: "bad-option",
		},
		{
			TestName: "success",
			Opts:     []dbkit.SelectOption{},
			SongRepoFn: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT id, title, artist, updated_at, created_at FROM songs`).
					WillReturnRows(sqlmock.
						NewRows([]string{"id", "title", "artist", "updated_at", "created_at"}).
						AddRow("1234", "some-title4", "some-artist4", now, now).
						AddRow("1235", "some-title5", "some-artist5", now, now),
					)
			},
			Expected: []*mysqldb.Song{
				&mysqldb.Song{ID: 1234, Title: "some-title4", Artist: "some-artist4", UpdatedAt: now, CreatedAt: now},
				&mysqldb.Song{ID: 1235, Title: "some-title5", Artist: "some-artist5", UpdatedAt: now, CreatedAt: now},
			},
		},
	}
	for _, tt := range testcases {
		t.Run(tt.TestName, func(t *testing.T) {
			repo, db := createSongRepo(tt.SongRepoFn)
			defer db.Close()

			songs, err := repo.Find(context.Background(), tt.Opts...)
			if tt.ExpectedErr != "" {
				require.EqualError(t, err, tt.ExpectedErr)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.Expected, songs)
			}
		})
	}
}

func TestSongRepoImpl_Delete(t *testing.T) {
	testcases := []struct {
		TestName    string
		Opt         dbkit.DeleteOption
		SongRepoFn  bookRepoFn
		ExpectedErr string
		Expected    int64
	}{
		{
			TestName:    "begin error",
			Opt:         dbkit.Equal("id", 666),
			ExpectedErr: "dbtxn: begin-error",
			SongRepoFn: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin().WillReturnError(errors.New("begin-error"))
			},
		},
		{
			TestName:    "delete error",
			Opt:         dbkit.Equal("id", 666),
			ExpectedErr: "delete-error",
			SongRepoFn: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM songs WHERE id = ?`)).
					WithArgs(666).
					WillReturnError(errors.New("delete-error"))
			},
		},
		{
			TestName: "bad delete option",
			Opt: dbkit.NewDeleteOption(func(b sq.DeleteBuilder) (sq.DeleteBuilder, error) {
				return b, errors.New("bad-option")
			}),
			SongRepoFn: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM songs WHERE id = ?`)).
					WithArgs(555).
					WillReturnResult(sqlmock.NewResult(0, 1))
			},
			ExpectedErr: "bad-option",
		},
		{
			TestName: "success",
			Opt:      dbkit.Equal("id", 555),
			SongRepoFn: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM songs WHERE id = ?`)).
					WithArgs(555).
					WillReturnResult(sqlmock.NewResult(0, 1))
			},
			Expected: 1,
		},
	}
	for _, tt := range testcases {
		t.Run(tt.TestName, func(t *testing.T) {
			repo, db := createSongRepo(tt.SongRepoFn)
			defer db.Close()

			ctx := context.Background()
			dbtxn.Begin(&ctx)

			affectedRow, err := repo.Delete(ctx, tt.Opt)
			if tt.ExpectedErr != "" {
				require.EqualError(t, err, tt.ExpectedErr)
				require.EqualError(t, dbtxn.Error(ctx), tt.ExpectedErr)
			} else {
				require.NoError(t, err)
				require.NoError(t, dbtxn.Error(ctx))
				require.Equal(t, tt.Expected, affectedRow)
			}
		})
	}
}
