package dbkit

import (
	"errors"

	sq "github.com/Masterminds/squirrel"
)

var (
	_ SelectOption = (*LikeOption)(nil)
)

type (

	// LikeOption for where condition
	LikeOption struct {
		column      string
		expectation interface{}
	}
)

//

// Like where condition
func Like(column string, expectation interface{}) *LikeOption {
	return &LikeOption{
		column:      column,
		expectation: expectation,
	}
}

// CompileSelect for like
func (l *LikeOption) CompileSelect(base sq.SelectBuilder) (sq.SelectBuilder, error) {
	if l.column == "" {
		return base, errors.New("Like column can't be empty")
	}
	return base.Where(sq.Like{l.column: l.expectation}), nil
}
