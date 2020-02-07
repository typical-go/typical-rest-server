package dbkit_test

import (
	"testing"

	sq "github.com/Masterminds/squirrel"
	"github.com/stretchr/testify/require"

	"github.com/typical-go/typical-rest-server/pkg/dbkit"
)

func TestSort(t *testing.T) {
	testcases := []struct {
		dbkit.FindOption
		builder       sq.SelectBuilder
		expectedError string
		expected      string
	}{
		{
			FindOption:    dbkit.Sort("", 0),
			builder:       sq.Select("name", "version").From("sometables"),
			expectedError: "Sort column can't be empty",
		},
		{
			FindOption: dbkit.Sort("name", dbkit.Asc),
			builder:    sq.Select("name", "version").From("sometables"),
			expected:   "SELECT name, version FROM sometables ORDER BY name ASC",
		},
		{
			FindOption: dbkit.Sort("other_col", dbkit.Desc),
			builder:    sq.Select("name", "version").From("sometables"),
			expected:   "SELECT name, version FROM sometables ORDER BY other_col DESC",
		},
	}

	for _, tt := range testcases {
		builder, err := tt.CompileQuery(tt.builder)
		if err != nil {
			require.EqualError(t, err, tt.expectedError)
		} else {
			query, _, _ := builder.ToSql()
			require.Equal(t, tt.expected, query)
		}
	}
}

func TestSortOrder(t *testing.T) {
	testcases := []struct {
		dbkit.OrderBy
		s string
	}{
		{OrderBy: dbkit.Asc, s: "ASC"},
		{OrderBy: dbkit.Desc, s: "DESC"},
		{OrderBy: 9999999, s: "ASC"},
	}
	for _, tt := range testcases {
		require.Equal(t, tt.s, tt.OrderBy.String())
	}

}
