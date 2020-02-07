package dbkit_test

import (
	"testing"

	sq "github.com/Masterminds/squirrel"
	"github.com/stretchr/testify/require"

	"github.com/typical-go/typical-rest-server/pkg/dbkit"
)

func TestPagination(t *testing.T) {
	testcases := []struct {
		dbkit.FindOption
		builder       sq.SelectBuilder
		expectedError string
		expected      string
	}{
		{
			FindOption:    dbkit.Pagination(0, 0),
			builder:       sq.Select("name", "version").From("sometables"),
			expectedError: "Limit can't be 0 or negative",
		},
		{
			FindOption: dbkit.Pagination(10, 100),
			builder:    sq.Select("name", "version").From("sometables"),
			expected:   "SELECT name, version FROM sometables LIMIT 100 OFFSET 10",
		},
		{
			FindOption: dbkit.PaginationWithRange(10, 100),
			builder:    sq.Select("name", "version").From("sometables"),
			expected:   "SELECT name, version FROM sometables LIMIT 91 OFFSET 10",
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
