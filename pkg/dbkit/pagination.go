package dbkit

import (
	"errors"

	sq "github.com/Masterminds/squirrel"
)

// Pagination param
type pagination struct {
	offset uint64
	limit  uint64
}

// Pagination find option
func Pagination(offset, limit uint64) FindOption {
	return &pagination{
		offset: offset,
		limit:  limit,
	}
}

// PaginationWithRange to setup pagination with start and end index
func PaginationWithRange(start, end uint64) FindOption {
	return Pagination(start, end-start+1)
}

// CompileQuery to compile select query for pagination
func (p *pagination) CompileQuery(base sq.SelectBuilder) (sq.SelectBuilder, error) {
	if p.limit < 1 {
		return base, errors.New("Limit can't be 0 or negative")
	}
	base = base.Offset(p.offset)
	if p.limit != 0 {
		base = base.Limit(p.limit)
	}
	return base, nil
}
