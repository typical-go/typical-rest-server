package dbkit

import (
	"errors"

	sq "github.com/Masterminds/squirrel"
)

// LikeOption for where condition
type LikeOption struct {
	column      string
	expectation interface{}
}

var _ SelectOption = (*LikeOption)(nil)

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
