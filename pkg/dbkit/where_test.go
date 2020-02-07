package dbkit_test

import (
	"testing"

	sq "github.com/Masterminds/squirrel"
	"github.com/stretchr/testify/require"
	"github.com/typical-go/typical-rest-server/pkg/dbkit"
)

func TestFindOption_Where(t *testing.T) {
	testcases := []struct {
		dbkit.WhereCondition
		builder         sq.SelectBuilder
		expectedError   string
		expected        string
		expectedSqlArgs []interface{}
	}{
		{
			WhereCondition: dbkit.Equal("", ""),
			builder:        sq.Select("name", "version").From("sometables"),
			expectedError:  "Filter column can't be empty",
		},
		{
			WhereCondition:  dbkit.Equal("name", "dummy-name"),
			builder:         sq.Select("name", "version").From("sometables"),
			expected:        "SELECT name, version FROM sometables WHERE name = ?",
			expectedSqlArgs: []interface{}{"dummy-name"},
		},
		{
			WhereCondition:  dbkit.Like("name", "%dum%"),
			builder:         sq.Select("name", "version").From("sometables"),
			expected:        "SELECT name, version FROM sometables WHERE name LIKE ?",
			expectedSqlArgs: []interface{}{"%dum%"},
		},
	}

	for _, tt := range testcases {
		builder, err := tt.CompileQuery(tt.builder)
		if err != nil {
			require.EqualError(t, err, tt.expectedError)
		} else {
			query, args, _ := builder.ToSql()
			require.Equal(t, tt.expected, query)
			require.Equal(t, tt.expectedSqlArgs, args)
		}
	}
}
