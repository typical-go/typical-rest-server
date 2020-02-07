package dbkit

import (
	"errors"
	"fmt"

	sq "github.com/Masterminds/squirrel"
)

// Sort param
type Sort struct {
	Column string
	Order  Order
}

type Order int

const (
	Asc Order = iota
	Desc
)

func (o Order) String() string {
	switch o {
	case Asc:
		return "ASC"
	case Desc:
		return "DESC"
	}
	return "ASC"
}

// NewSort return new instance of Sort
func NewSort(column string, order Order) *Sort {
	return &Sort{
		Column: column,
		Order:  order,
	}
}

// CompileQuery to compile select query for sorting
func (s *Sort) CompileQuery(base sq.SelectBuilder) (sq.SelectBuilder, error) {
	if s.Column == "" {
		return base, errors.New("Sort column can't be empty")
	}

	base = base.OrderBy(fmt.Sprintf("%s %s", s.Column, s.Order))
	return base, nil
}
