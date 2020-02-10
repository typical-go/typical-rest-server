package dbkit

import (
	"errors"

	sq "github.com/Masterminds/squirrel"
)

// WhereCondition is find option with WHERE condition
type WhereCondition FindOption

//
// Filter
//

type equal struct {
	column    string
	condition string
}

// Equal where condition
func Equal(column string, cond string) WhereCondition {
	return &equal{
		column:    column,
		condition: cond,
	}
}

// CompileQuery to compile select query for filtering
func (f *equal) CompileQuery(base sq.SelectBuilder) (sq.SelectBuilder, error) {
	if f.column == "" {
		return base, errors.New("Filter column can't be empty")
	}
	return base.Where(sq.Eq{f.column: f.condition}), nil
}

//
// Like
//

type like struct {
	column    string
	condition string
}

// Like where condition
func Like(column, condition string) WhereCondition {
	return &like{
		column:    column,
		condition: condition,
	}
}

func (l *like) CompileQuery(base sq.SelectBuilder) (sq.SelectBuilder, error) {
	if l.column == "" {
		return base, errors.New("Like column can't be empty")
	}
	return base.Where(sq.Like{l.column: l.condition}), nil
}
