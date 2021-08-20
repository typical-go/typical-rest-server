package repo_test

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"
	"github.com/typical-go/typical-rest-server/internal/app/repo"
	"github.com/typical-go/typical-rest-server/pkg/sqkit"
)

func TestBookRepoImpl_Count(t *testing.T) {
	testCases := []struct {
		TestName       string
		OnMockDatabase func(mock sqlmock.Sqlmock)
		Opts           []sqkit.SelectOption
		Expected       []*repo.Book
		ExpectedErr    string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range testCases {
		t.Run(tt.TestName, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			require.NoError(t, err)
			defer db.Close()

			if tt.OnMockDatabase != nil {
				tt.OnMockDatabase(mock)
			}

			repo := &repo.BookRepoImpl{DB: db}

			books, err := repo.Count(context.Background(), tt.Opts...)
			if tt.ExpectedErr != "" {
				require.EqualError(t, err, tt.ExpectedErr)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.Expected, books)
			}
		})
	}
}

func TestBookRepoImpl_Find(t *testing.T) {
	testCases := []struct {
		TestName       string
		OnMockDatabase func(mock sqlmock.Sqlmock)
		Opts           []sqkit.SelectOption
		Expected       []*repo.Book
		ExpectedErr    string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range testCases {
		t.Run(tt.TestName, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			require.NoError(t, err)
			defer db.Close()

			if tt.OnMockDatabase != nil {
				tt.OnMockDatabase(mock)
			}

			repo := &repo.BookRepoImpl{DB: db}

			books, err := repo.Find(context.Background(), tt.Opts...)
			if tt.ExpectedErr != "" {
				require.EqualError(t, err, tt.ExpectedErr)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.Expected, books)
			}
		})
	}
}

func TestBookRepoImpl_Insert(t *testing.T) {
	testCases := []struct {
		TestName       string
		OnMockDatabase func(mock sqlmock.Sqlmock)
		Book           *repo.Book
		Expected       int64
		ExpectedErr    string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range testCases {
		t.Run(tt.TestName, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			require.NoError(t, err)
			defer db.Close()

			if tt.OnMockDatabase != nil {
				tt.OnMockDatabase(mock)
			}

			repo := &repo.BookRepoImpl{DB: db}

			id, err := repo.Insert(context.Background(), tt.Book)
			if tt.ExpectedErr != "" {
				require.EqualError(t, err, tt.ExpectedErr)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.Expected, id)
			}
		})
	}
}

func TestBookRepoImpl_BulkInsert(t *testing.T) {
	testCases := []struct {
		TestName       string
		OnMockDatabase func(mock sqlmock.Sqlmock)
		Books          []*repo.Book
		Expected       int64
		ExpectedErr    string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range testCases {
		t.Run(tt.TestName, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			require.NoError(t, err)
			defer db.Close()

			if tt.OnMockDatabase != nil {
				tt.OnMockDatabase(mock)
			}

			repo := &repo.BookRepoImpl{DB: db}

			affectedRows, err := repo.BulkInsert(context.Background(), tt.Books...)
			if tt.ExpectedErr != "" {
				require.EqualError(t, err, tt.ExpectedErr)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.Expected, affectedRows)
			}
		})
	}
}

func TestBookRepoImpl_Delete(t *testing.T) {
	// 		Delete(context.Context, sqkit.DeleteOption) (int64, error)
	testCases := []struct {
		TestName       string
		OnMockDatabase func(mock sqlmock.Sqlmock)
		Opts           []sqkit.DeleteOption
		Expected       int64
		ExpectedErr    string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range testCases {
		t.Run(tt.TestName, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			require.NoError(t, err)
			defer db.Close()

			if tt.OnMockDatabase != nil {
				tt.OnMockDatabase(mock)
			}

			repo := &repo.BookRepoImpl{DB: db}

			affectedRows, err := repo.Delete(context.Background(), tt.Opts...)
			if tt.ExpectedErr != "" {
				require.EqualError(t, err, tt.ExpectedErr)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.Expected, affectedRows)
			}
		})
	}
}

func TestBookRepoImpl_Update(t *testing.T) {
	testCases := []struct {
		TestName       string
		OnMockDatabase func(mock sqlmock.Sqlmock)
		Book           *repo.Book
		Opts           []sqkit.UpdateOption
		Expected       []*repo.Book
		ExpectedErr    string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range testCases {
		t.Run(tt.TestName, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			require.NoError(t, err)
			defer db.Close()

			if tt.OnMockDatabase != nil {
				tt.OnMockDatabase(mock)
			}

			repo := &repo.BookRepoImpl{DB: db}

			books, err := repo.Update(context.Background(), tt.Book, tt.Opts...)
			if tt.ExpectedErr != "" {
				require.EqualError(t, err, tt.ExpectedErr)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.Expected, books)
			}
		})
	}
}

func TestBookRepoImpl_Patch(t *testing.T) {
	testCases := []struct {
		TestName       string
		OnMockDatabase func(mock sqlmock.Sqlmock)
		Book           *repo.Book
		Opts           []sqkit.UpdateOption
		Expected       []*repo.Book
		ExpectedErr    string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range testCases {
		t.Run(tt.TestName, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			require.NoError(t, err)
			defer db.Close()

			if tt.OnMockDatabase != nil {
				tt.OnMockDatabase(mock)
			}

			repo := &repo.BookRepoImpl{DB: db}

			books, err := repo.Patch(context.Background(), tt.Book, tt.Opts...)
			if tt.ExpectedErr != "" {
				require.EqualError(t, err, tt.ExpectedErr)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.Expected, books)
			}
		})
	}
}
