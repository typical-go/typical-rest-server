package dbkit

import (
	"errors"
	"fmt"

	sq "github.com/Masterminds/squirrel"
)

type sort struct {
	column string
	order  Order
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

// Sort is find option to sort by column and order
func Sort(column string, order Order) FindOption {
	return &sort{
		column: column,
		order:  order,
	}
}

// CompileQuery to compile select query for sorting
func (s *sort) CompileQuery(base sq.SelectBuilder) (sq.SelectBuilder, error) {
	if s.column == "" {
		return base, errors.New("Sort column can't be empty")
	}
	base = base.OrderBy(fmt.Sprintf("%s %s", s.column, s.order))
	return base, nil
}
