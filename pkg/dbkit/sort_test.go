package dbkit_test

import (
	"testing"

	sq "github.com/Masterminds/squirrel"
	"github.com/stretchr/testify/require"

	"github.com/typical-go/typical-rest-server/pkg/dbkit"
)

func TestSort(t *testing.T) {
	testcases := []struct {
		*dbkit.Sort
		builder       sq.SelectBuilder
		expectedError string
		expected      string
	}{
		{
			Sort:          &dbkit.Sort{},
			builder:       sq.Select("name", "version").From("sometables"),
			expectedError: "Sort column can't be empty",
		},
		{
			Sort:     dbkit.NewSort("name", dbkit.Asc),
			builder:  sq.Select("name", "version").From("sometables"),
			expected: "SELECT name, version FROM sometables ORDER BY name ASC",
		},
		{
			Sort:     dbkit.NewSort("other_col", dbkit.Desc),
			builder:  sq.Select("name", "version").From("sometables"),
			expected: "SELECT name, version FROM sometables ORDER BY other_col DESC",
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
