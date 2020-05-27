package dbkit_test

import (
	"testing"

	sq "github.com/Masterminds/squirrel"
	"github.com/stretchr/testify/require"

	"github.com/typical-go/typical-rest-server/pkg/dbkit"
)

func TestSort(t *testing.T) {
	testcases := []dbkit.SelectTestCase{
		{
			SelectOption: dbkit.Sort("", 0),
			Builder:      sq.Select("name", "version").From("sometables"),
			ExpectedErr:  "Sort column can't be empty",
		},
		{
			SelectOption: dbkit.Sort("name", dbkit.Asc),
			Builder:      sq.Select("name", "version").From("sometables"),
			Expected:     "SELECT name, version FROM sometables ORDER BY name ASC",
		},
		{
			SelectOption: dbkit.Sort("other_col", dbkit.Desc),
			Builder:      sq.Select("name", "version").From("sometables"),
			Expected:     "SELECT name, version FROM sometables ORDER BY other_col DESC",
		},
	}

	for _, tt := range testcases {
		tt.Execute(t)
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
