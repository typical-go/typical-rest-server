package dbkit

import (
	"errors"

	sq "github.com/Masterminds/squirrel"
)

var (
	_ SelectOption = (*EqualOption)(nil)
	_ UpdateOption = (*EqualOption)(nil)
	_ DeleteOption = (*EqualOption)(nil)
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
		return base, errors.New("equal: column is missing")
	}
	return base.Where(sq.Eq{f.column: f.expectation}), nil
}

// CompileUpdate to compile update query for filtering
func (f *EqualOption) CompileUpdate(base sq.UpdateBuilder) (sq.UpdateBuilder, error) {
	if f.column == "" {
		return base, errors.New("equal: column is missing")
	}
	return base.Where(sq.Eq{f.column: f.expectation}), nil
}

// CompileDelete to compile delete query for filtering
func (f *EqualOption) CompileDelete(base sq.DeleteBuilder) (sq.DeleteBuilder, error) {
	if f.column == "" {
		return base, errors.New("equal: column is missing")
	}
	return base.Where(sq.Eq{f.column: f.expectation}), nil
}
