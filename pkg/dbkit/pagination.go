package dbkit

import (
	"errors"

	sq "github.com/Masterminds/squirrel"
)

// Pagination param
type Pagination struct {
	Offset uint64
	Limit  uint64
}

// NewPagination return new instance of Pagination
func NewPagination(offset, limit uint64) *Pagination {
	return &Pagination{
		Offset: offset,
		Limit:  limit,
	}

}

// CreatePaginationWithRange to create pagination
func CreatePaginationWithRange(start, end uint64) *Pagination {
	return NewPagination(start, end-start+1)
}

// CompileQuery to compile select query for pagination
func (p *Pagination) CompileQuery(base sq.SelectBuilder) (sq.SelectBuilder, error) {
	if p.Limit < 1 {
		return base, errors.New("Limit can't be 0 or negative")
	}
	base = base.Offset(p.Offset)
	if p.Limit != 0 {
		base = base.Limit(p.Limit)
	}
	return base, nil
}
