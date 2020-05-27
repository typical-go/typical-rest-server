package dbkit_test

import (
	"testing"

	sq "github.com/Masterminds/squirrel"

	"github.com/typical-go/typical-rest-server/pkg/dbkit"
)

func TestPagination(t *testing.T) {
	testcases := []dbkit.SelectTestCase{
		{
			SelectOption: dbkit.Pagination(0, 0),
			Builder:      sq.Select("name", "version").From("sometables"),
			ExpectedErr:  "Limit can't be 0 or negative",
		},
		{
			SelectOption: dbkit.Pagination(10, 100),
			Builder:      sq.Select("name", "version").From("sometables"),
			Expected:     "SELECT name, version FROM sometables LIMIT 100 OFFSET 10",
		},
		{
			SelectOption: dbkit.PaginationWithRange(10, 100),
			Builder:      sq.Select("name", "version").From("sometables"),
			Expected:     "SELECT name, version FROM sometables LIMIT 91 OFFSET 10",
		},
	}

	for _, tt := range testcases {
		tt.Execute(t)
	}
}
