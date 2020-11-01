package sqkit_test

import (
	"testing"

	sq "github.com/Masterminds/squirrel"
	"github.com/stretchr/testify/require"
	"github.com/typical-go/typical-rest-server/pkg/sqkit"
)

func TestEq_CompileSelect(t *testing.T) {
	testcases := []struct {
		testName string
		sqkit.Eq
		base          sq.SelectBuilder
		expectedQuery string
		expectedArgs  []interface{}
	}{
		{
			base:          sq.Select("name", "version").From("some-table"),
			expectedQuery: "SELECT name, version FROM some-table",
		},
		{
			Eq:            sqkit.Eq{"name": "dummy-name"},
			base:          sq.Select("name", "version").From("some-table"),
			expectedQuery: "SELECT name, version FROM some-table WHERE name = ?",
			expectedArgs:  []interface{}{"dummy-name"},
		},
	}

	for _, tt := range testcases {
		t.Run(tt.testName, func(t *testing.T) {
			query, args, err := tt.CompileSelect(tt.base).ToSql()
			require.NoError(t, err)
			require.Equal(t, tt.expectedQuery, query)
			require.Equal(t, tt.expectedArgs, args)
		})
	}
}

func TestEq_CompileUpdate(t *testing.T) {
	testcases := []struct {
		testName string
		sqkit.Eq
		base          sq.UpdateBuilder
		expectedQuery string
		expectedArgs  []interface{}
	}{
		{
			Eq:            sqkit.Eq{"name": "dummy-name"},
			base:          sq.Update("some-table").Set("column", "column-value"),
			expectedQuery: "UPDATE some-table SET column = ? WHERE name = ?",
			expectedArgs:  []interface{}{"column-value", "dummy-name"},
		},
	}

	for _, tt := range testcases {
		t.Run(tt.testName, func(t *testing.T) {
			query, args, _ := tt.CompileUpdate(tt.base).ToSql()
			require.Equal(t, tt.expectedQuery, query)
			require.Equal(t, tt.expectedArgs, args)
		})
	}

}

func TestEq_Delete(t *testing.T) {
	testcases := []struct {
		testName string
		sqkit.Eq
		base          sq.DeleteBuilder
		expectedQuery string
		expectedArgs  []interface{}
	}{
		{
			Eq:            sqkit.Eq{"name": "dummy-name"},
			base:          sq.Delete("some-table"),
			expectedQuery: "DELETE FROM some-table WHERE name = ?",
			expectedArgs:  []interface{}{"dummy-name"},
		},
	}

	for _, tt := range testcases {
		t.Run(tt.testName, func(t *testing.T) {
			query, args, _ := tt.CompileDelete(tt.base).ToSql()
			require.Equal(t, tt.expectedQuery, query)
			require.Equal(t, tt.expectedArgs, args)
		})
	}
}
