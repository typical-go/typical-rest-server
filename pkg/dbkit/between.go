package dbkit

import (
	"errors"

	sq "github.com/Masterminds/squirrel"
)

// BetweenOption is a placeholder for value used in a BETWEEN clause
type BetweenOption struct {
	column   string
	minValue interface{}
	maxValue interface{}
}

var _ SelectOption = (*BetweenOption)(nil)
var _ UpdateOption = (*BetweenOption)(nil)
var _ DeleteOption = (*BetweenOption)(nil)

// Between creates a new BetweenOption condition
func Between(column string, minValue, maxValue interface{}) *BetweenOption {
	return &BetweenOption{
		column:   column,
		minValue: minValue,
		maxValue: maxValue,
	}
}

// CompileSelect compile select query for filtering
func (f *BetweenOption) CompileSelect(base sq.SelectBuilder) (sq.SelectBuilder, error) {
	if f.column == "" {
		return base, errors.New("between: column is missing")
	}
	return base.Where(sq.And{
		sq.GtOrEq{f.column: f.minValue},
		sq.LtOrEq{f.column: f.maxValue},
	}), nil
}

// CompileUpdate compile update query for filtering
func (f *BetweenOption) CompileUpdate(base sq.UpdateBuilder) (sq.UpdateBuilder, error) {
	if f.column == "" {
		return base, errors.New("between: column is missing")
	}
	return base.Where(sq.And{
		sq.GtOrEq{f.column: f.minValue},
		sq.LtOrEq{f.column: f.maxValue},
	}), nil
}

// CompileDelete compile delete query for filtering
func (f *BetweenOption) CompileDelete(base sq.DeleteBuilder) (sq.DeleteBuilder, error) {
	if f.column == "" {
		return base, errors.New("between: column is missing")
	}
	return base.Where(sq.And{
		sq.GtOrEq{f.column: f.minValue},
		sq.LtOrEq{f.column: f.maxValue},
	}), nil
}
