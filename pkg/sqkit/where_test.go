package sqkit_test

import (
	"testing"

	sq "github.com/Masterminds/squirrel"
	"github.com/stretchr/testify/require"
	"github.com/typical-go/typical-rest-server/pkg/sqkit"
)

func TestWhere_CompileSelect(t *testing.T) {
	testcases := []struct {
		testName string
		sqkit.Where
		base          sq.SelectBuilder
		expectedQuery string
		expectedArgs  []interface{}
	}{
		{
			base:          sq.Select("name", "version").From("some-table"),
			expectedQuery: "SELECT name, version FROM some-table",
		},
		{
			Where: sqkit.Where{
				sq.Eq{"name": "dummy-name"},
				sq.GtOrEq{"version": 1},
			},
			base:          sq.Select("name", "version").From("some-table"),
			expectedQuery: "SELECT name, version FROM some-table WHERE name = ? AND version >= ?",
			expectedArgs:  []interface{}{"dummy-name", 1},
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

func TestWhere_CompileUpdate(t *testing.T) {
	testcases := []struct {
		testName string
		sqkit.Where
		base          sq.UpdateBuilder
		expectedQuery string
		expectedArgs  []interface{}
	}{
		{
			Where: sqkit.Where{
				sq.Eq{"name": "dummy-name"},
				sq.LtOrEq{"version": 2},
			},
			base:          sq.Update("some-table").Set("column", "column-value"),
			expectedQuery: "UPDATE some-table SET column = ? WHERE name = ? AND version <= ?",
			expectedArgs:  []interface{}{"column-value", "dummy-name", 2},
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

func TestWhere_Delete(t *testing.T) {
	testcases := []struct {
		testName string
		sqkit.Where
		base          sq.DeleteBuilder
		expectedQuery string
		expectedArgs  []interface{}
	}{
		{
			Where: sqkit.Where{
				sq.Eq{"name": "dummy-name"},
				sq.Lt{"version": 3},
			},
			base:          sq.Delete("some-table"),
			expectedQuery: "DELETE FROM some-table WHERE name = ? AND version < ?",
			expectedArgs:  []interface{}{"dummy-name", 3},
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
