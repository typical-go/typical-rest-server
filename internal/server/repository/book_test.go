package repository_test

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"testing"
	"time"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"
	"github.com/typical-go/typical-rest-server/internal/server/repository"
	"github.com/typical-go/typical-rest-server/pkg/dbkit"
)

type (
	bookFindTestCase struct {
		testName    string
		opts        []dbkit.SelectOption
		expected    []*repository.Book
		expectedErr string
		mockFn      func(sqlmock.Sqlmock)
	}

	bookCreateTestCase struct {
		testName    string
		book        *repository.Book
		expected    int64
		expectedErr string
		mockFn      func(sqlmock.Sqlmock)
	}

	bookUpdateTestCase struct {
		testName    string
		book        *repository.Book
		opt         dbkit.UpdateOption
		expectedErr string
		mockFn      func(sqlmock.Sqlmock)
	}

	bookDeleteTestCase struct {
		testName    string
		opt         dbkit.DeleteOption
		expectedErr string
		mockFn      func(sqlmock.Sqlmock)
	}
)

func TestBookRepoImpl_Create(t *testing.T) {
	testcases := []bookCreateTestCase{
		{
			book: &repository.Book{
				Title:  "some-title",
				Author: "some-author",
			},
			expectedErr: "some-error",
			mockFn: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO books (title,author,created_at,updated_at) VALUES ($1,$2,$3,$4) RETURNING "id"`)).
					WithArgs("some-title", "some-author", sqlmock.AnyArg(), sqlmock.AnyArg()).
					WillReturnError(errors.New("some-error"))
			},
		},
		{
			book: &repository.Book{
				Title:  "some-title",
				Author: "some-author",
			},
			expected: 999,
			mockFn: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO books (title,author,created_at,updated_at) VALUES ($1,$2,$3,$4) RETURNING "id"`)).
					WithArgs("some-title", "some-author", sqlmock.AnyArg(), sqlmock.AnyArg()).
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(999))
			},
		},
	}

	for _, tt := range testcases {
		t.Run(tt.testName, func(t *testing.T) {
			db, mock, _ := sqlmock.New()
			defer db.Close()
			if tt.mockFn != nil {
				tt.mockFn(mock)
			}
			repo := &repository.BookRepoImpl{DB: db}
			id, err := repo.Create(context.Background(), tt.book)
			if tt.expectedErr != "" {
				require.EqualError(t, err, tt.expectedErr)
				return
			}
			require.NoError(t, err)
			require.Equal(t, tt.expected, id)
		})
	}
}

func TestBookRepoImpl_Update(t *testing.T) {
	testcases := []bookUpdateTestCase{
		{
			book: &repository.Book{
				Title:  "new-title",
				Author: "new-author",
			},
			opt:         dbkit.Equal(repository.BookCols.ID, 888),
			expectedErr: "some-update-error",
			mockFn: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(regexp.QuoteMeta(`UPDATE books SET title = $1, author = $2, updated_at = $3 WHERE id = $4`)).
					WithArgs("new-title", "new-author", sqlmock.AnyArg(), 888).
					WillReturnError(fmt.Errorf("some-update-error"))
			},
		},
		{
			book: &repository.Book{
				Title:  "new-title",
				Author: "new-author",
			},
			opt: dbkit.Equal(repository.BookCols.ID, 888),
			mockFn: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(regexp.QuoteMeta(`UPDATE books SET title = $1, author = $2, updated_at = $3 WHERE id = $4`)).
					WithArgs("new-title", "new-author", sqlmock.AnyArg(), 888).
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
		},
		{
			testName: "empty author",
			book: &repository.Book{
				Title: "new-title",
			},
			opt: dbkit.Equal(repository.BookCols.ID, 888),
			mockFn: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(regexp.QuoteMeta(`UPDATE books SET title = $1, author = $2, updated_at = $3 WHERE id = $4`)).
					WithArgs("new-title", "", sqlmock.AnyArg(), 888).
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
		},
		{
			testName: "empty title",
			book: &repository.Book{
				Author: "new-author",
			},
			opt: dbkit.Equal(repository.BookCols.ID, 888),
			mockFn: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(regexp.QuoteMeta(`UPDATE books SET title = $1, author = $2, updated_at = $3 WHERE id = $4`)).
					WithArgs("", "new-author", sqlmock.AnyArg(), 888).
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
		},
	}
	for _, tt := range testcases {
		t.Run(tt.testName, func(t *testing.T) {
			db, mock, _ := sqlmock.New()
			defer db.Close()
			if tt.mockFn != nil {
				tt.mockFn(mock)
			}
			repo := &repository.BookRepoImpl{DB: db}
			err := repo.Update(context.Background(), tt.book, tt.opt)
			if tt.expectedErr != "" {
				require.EqualError(t, err, tt.expectedErr)
				return
			}
			require.NoError(t, err)
		})
	}
}

func TestBookRepoImpl_Patch(t *testing.T) {
	testcases := []bookUpdateTestCase{
		{
			book: &repository.Book{
				Title:  "new-title",
				Author: "new-author",
			},
			opt:         dbkit.Equal(repository.BookCols.ID, 888),
			expectedErr: "some-update-error",
			mockFn: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(regexp.QuoteMeta(`UPDATE books SET title = $1, author = $2, updated_at = $3 WHERE id = $4`)).
					WithArgs("new-title", "new-author", sqlmock.AnyArg(), 888).
					WillReturnError(fmt.Errorf("some-update-error"))
			},
		},
		{
			book: &repository.Book{
				Title:  "new-title",
				Author: "new-author",
			},
			opt: dbkit.Equal(repository.BookCols.ID, 888),
			mockFn: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(regexp.QuoteMeta(`UPDATE books SET title = $1, author = $2, updated_at = $3 WHERE id = $4`)).
					WithArgs("new-title", "new-author", sqlmock.AnyArg(), 888).
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
		},
		{
			testName: "empty author",
			book: &repository.Book{
				Title: "new-title",
			},
			opt: dbkit.Equal(repository.BookCols.ID, 888),
			mockFn: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(regexp.QuoteMeta(`UPDATE books SET title = $1, updated_at = $2 WHERE id = $3`)).
					WithArgs("new-title", sqlmock.AnyArg(), 888).
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
		},
		{
			testName: "empty title",
			book: &repository.Book{
				Author: "new-author",
			},
			opt: dbkit.Equal(repository.BookCols.ID, 888),
			mockFn: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(regexp.QuoteMeta(`UPDATE books SET author = $1, updated_at = $2 WHERE id = $3`)).
					WithArgs("new-author", sqlmock.AnyArg(), 888).
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
		},
	}
	for _, tt := range testcases {
		t.Run(tt.testName, func(t *testing.T) {
			db, mock, _ := sqlmock.New()
			defer db.Close()
			if tt.mockFn != nil {
				tt.mockFn(mock)
			}
			repo := &repository.BookRepoImpl{DB: db}
			err := repo.Patch(context.Background(), tt.book, tt.opt)
			if tt.expectedErr != "" {
				require.EqualError(t, err, tt.expectedErr)
				return
			}
			require.NoError(t, err)
		})
	}
}

func TestBookRepoImpl_Find(t *testing.T) {
	now := time.Now()
	testcases := []bookFindTestCase{
		{

			opts:        []dbkit.SelectOption{},
			expected:    []*repository.Book{},
			expectedErr: "some-error",
			mockFn: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT id, title, author, updated_at, created_at FROM books`).
					WillReturnError(errors.New("some-error"))
			},
		},
		{

			opts: []dbkit.SelectOption{},
			expected: []*repository.Book{
				&repository.Book{ID: 1234, Title: "some-title4", Author: "some-author4", UpdatedAt: now, CreatedAt: now},
				&repository.Book{ID: 1235, Title: "some-title5", Author: "some-author5", UpdatedAt: now, CreatedAt: now},
			},
			mockFn: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "title", "author", "updated_at", "created_at"})
				rows.AddRow("1234", "some-title4", "some-author4", now, now)
				rows.AddRow("1235", "some-title5", "some-author5", now, now)
				mock.ExpectQuery(`SELECT id, title, author, updated_at, created_at FROM books`).
					WillReturnRows(rows)
			},
		},
	}
	for _, tt := range testcases {
		t.Run(tt.testName, func(t *testing.T) {
			db, mock, _ := sqlmock.New()
			defer db.Close()
			if tt.mockFn != nil {
				tt.mockFn(mock)
			}
			repo := &repository.BookRepoImpl{DB: db}
			books, err := repo.Find(context.Background(), tt.opts...)
			if tt.expectedErr != "" {
				require.EqualError(t, err, tt.expectedErr)
				return
			}
			require.NoError(t, err)
			require.Equal(t, tt.expected, books)
		})
	}
}

func TestBookRepoImpl_Delete(t *testing.T) {
	testcases := []bookDeleteTestCase{
		{
			opt:         dbkit.Equal("id", 666),
			expectedErr: "some-error",
			mockFn: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM books WHERE id = $1`)).
					WithArgs(666).
					WillReturnError(fmt.Errorf("some-error"))
			},
		},
		{
			opt: dbkit.Equal("id", 555),
			mockFn: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM books WHERE id = $1`)).
					WithArgs(555).
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
		},
	}
	for _, tt := range testcases {
		t.Run(tt.testName, func(t *testing.T) {
			db, mock, _ := sqlmock.New()
			defer db.Close()
			if tt.mockFn != nil {
				tt.mockFn(mock)
			}
			repo := &repository.BookRepoImpl{DB: db}
			err := repo.Delete(context.Background(), tt.opt)
			if tt.expectedErr != "" {
				require.EqualError(t, err, tt.expectedErr)
				return
			}
			require.NoError(t, err)
		})
	}
}
