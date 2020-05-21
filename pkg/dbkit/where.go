package dbkit

import (
	"errors"

	sq "github.com/Masterminds/squirrel"
)

var (
	_ SelectOption = (*EqualOption)(nil)
	_ SelectOption = (*LikeOption)(nil)
)

type (
	// EqualOption for where condition
	EqualOption struct {
		column      string
		expectation interface{}
	}

	// LikeOption for where condition
	LikeOption struct {
		column      string
		expectation interface{}
	}
)

//
// Filter
//

// Equal where condition
func Equal(column string, expectation interface{}) *EqualOption {
	return &EqualOption{
		column:      column,
		expectation: expectation,
	}
}

// CompileQuery to compile select query for filtering
func (f *EqualOption) CompileQuery(base sq.SelectBuilder) (sq.SelectBuilder, error) {
	if f.column == "" {
		return base, errors.New("Filter column can't be empty")
	}
	return base.Where(sq.Eq{f.column: f.expectation}), nil
}

//
// Like
//

// Like where condition
func Like(column string, expectation interface{}) *LikeOption {
	return &LikeOption{
		column:      column,
		expectation: expectation,
	}
}

// CompileQuery for like
func (l *LikeOption) CompileQuery(base sq.SelectBuilder) (sq.SelectBuilder, error) {
	if l.column == "" {
		return base, errors.New("Like column can't be empty")
	}
	return base.Where(sq.Like{l.column: l.expectation}), nil
}
