package dbkit

import (
	"errors"
	"fmt"

	sq "github.com/Masterminds/squirrel"
)

const (
	// Asc for ascending
	Asc OrderBy = iota

	// Desc for descending
	Desc
)

var (
	_ SelectOption = (*SortOption)(nil)
)

type (
	// SortOption for select
	SortOption struct {
		column  string
		orderBy OrderBy
	}

	// OrderBy is type of order by
	OrderBy int
)

//
// OrderBy
//

func (o OrderBy) String() string {
	switch o {
	case Asc:
		return "ASC"
	case Desc:
		return "DESC"
	}
	return "ASC"
}

//
// Sort
//

// Sort is find option to sort by column and order
func Sort(column string, orderBy OrderBy) *SortOption {
	return &SortOption{
		column:  column,
		orderBy: orderBy,
	}
}

// CompileSelect to compile select query for sorting
func (s *SortOption) CompileSelect(base sq.SelectBuilder) (sq.SelectBuilder, error) {
	if s.column == "" {
		return base, errors.New("Sort column can't be empty")
	}
	base = base.OrderBy(fmt.Sprintf("%s %s", s.column, s.orderBy))
	return base, nil
}
