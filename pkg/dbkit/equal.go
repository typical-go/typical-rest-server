package dbkit

import (
	"errors"

	sq "github.com/Masterminds/squirrel"
)

var (
	_ SelectOption = (*EqualOption)(nil)
)

type (
	// EqualOption for where condition
	EqualOption struct {
		column      string
		expectation interface{}
	}
)

// Equal where condition
func Equal(column string, expectation interface{}) *EqualOption {
	return &EqualOption{
		column:      column,
		expectation: expectation,
	}
}

// CompileSelect to compile select query for filtering
func (f *EqualOption) CompileSelect(base sq.SelectBuilder) (sq.SelectBuilder, error) {
	if f.column == "" {
		return base, errors.New("Filter column can't be empty")
	}
	return base.Where(sq.Eq{f.column: f.expectation}), nil
}
