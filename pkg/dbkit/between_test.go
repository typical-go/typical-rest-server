package dbkit_test

import (
	"testing"

	sq "github.com/Masterminds/squirrel"
	"github.com/stretchr/testify/assert"
	"github.com/typical-go/typical-rest-server/pkg/dbkit"
)

func TestBetween_Implementation(t *testing.T) {
	t.Run("Implements dbkit.SelectOption", func(t *testing.T) {
		assert.Implements(t, (*dbkit.SelectOption)(nil), new(dbkit.BetweenOption))
	})
	t.Run("Implements dbkit.UpdateOption", func(t *testing.T) {
		assert.Implements(t, (*dbkit.UpdateOption)(nil), new(dbkit.BetweenOption))
	})
	t.Run("Implements dbkit.DeleteOption", func(t *testing.T) {
		assert.Implements(t, (*dbkit.DeleteOption)(nil), new(dbkit.BetweenOption))
	})
}

func TestBetweenOption_Select(t *testing.T) {
	testCases := []dbkit.SelectTestCase{
		{
			SelectOption: dbkit.Between("", "", ""),
			Builder:      sq.Select("name", "version").From("some-table"),
			ExpectedErr:  "between: column is missing",
		},
		{
			SelectOption: dbkit.Between("date", "2020-01-01", "2020-01-31"),
			Builder:      sq.Select("name", "version").From("some-table"),
			Expected:     "SELECT name, version FROM some-table WHERE (date >= ? AND date <= ?)",
			ExpectedArgs: []interface{}{"2020-01-01", "2020-01-31"},
		},
	}
	for _, tt := range testCases {
		tt.Execute(t)
	}
}

func TestBetweenOption_Update(t *testing.T) {
	testCases := []dbkit.UpdateTestCase{
		{
			UpdateOption: dbkit.Between("", "", ""),
			Builder:      sq.Update("some-table"),
			ExpectedErr:  "between: column is missing",
		},
		{
			UpdateOption: dbkit.Between("date", "2020-01-01", "2020-01-31"),
			Builder:      sq.Update("some-table").Set("column", "column-value"),
			Expected:     "UPDATE some-table SET column = ? WHERE (date >= ? AND date <= ?)",
			ExpectedArgs: []interface{}{"column-value", "2020-01-01", "2020-01-31"},
		},
	}
	for _, tt := range testCases {
		tt.Execute(t)
	}
}

func TestBetweenOption_Delete(t *testing.T) {
	testCases := []dbkit.DeleteTestCase{
		{
			DeleteOption: dbkit.Between("", "", ""),
			Builder:      sq.Delete("some-table"),
			ExpectedErr:  "between: column is missing",
		},
		{
			DeleteOption: dbkit.Between("date", "2020-01-01", "2020-01-31"),
			Builder:      sq.Delete("some-table"),
			Expected:     "DELETE FROM some-table WHERE (date >= ? AND date <= ?)",
			ExpectedArgs: []interface{}{"2020-01-01", "2020-01-31"},
		},
	}
	for _, tt := range testCases {
		tt.Execute(t)
	}
}
