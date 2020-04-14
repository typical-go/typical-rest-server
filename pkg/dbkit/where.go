package dbkit

import (
	"errors"
	"fmt"

	sq "github.com/Masterminds/squirrel"
)

// WhereCondition is find option with WHERE condition
type WhereCondition FindOption

//
// Filter
//

type equal struct {
	column      string
	expectation interface{}
}

// Equal where condition
func Equal(column string, expectation interface{}) WhereCondition {
	return &equal{
		column:      column,
		expectation: expectation,
	}
}

// CompileQuery to compile select query for filtering
func (f *equal) CompileQuery(base sq.SelectBuilder) (sq.SelectBuilder, error) {
	if f.column == "" {
		return base, errors.New("Filter column can't be empty")
	}
	return base.Where(sq.Eq{f.column: f.expectation}), nil
}

func (f equal) String() string {
	return fmt.Sprintf("equal:%s:%v", f.column, f.expectation)
}

//
// Like
//

type like struct {
	column      string
	expectation interface{}
}

// Like where condition
func Like(column string, expectation interface{}) WhereCondition {
	return &like{
		column:      column,
		expectation: expectation,
	}
}

func (l *like) CompileQuery(base sq.SelectBuilder) (sq.SelectBuilder, error) {
	if l.column == "" {
		return base, errors.New("Like column can't be empty")
	}
	return base.Where(sq.Like{l.column: l.expectation}), nil
}

func (l like) String() string {
	return fmt.Sprintf("like:%s:%v", l.column, l.expectation)
}
