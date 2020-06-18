package dbkit

import (
	"errors"

	sq "github.com/Masterminds/squirrel"
)

// PaginationOption param
type PaginationOption struct {
	offset uint64
	limit  uint64
}

var _ SelectOption = (*PaginationOption)(nil)

// Pagination find option
func Pagination(offset, limit uint64) *PaginationOption {
	return &PaginationOption{
		offset: offset,
		limit:  limit,
	}
}

// PaginationWithRange to setup pagination with start and end index
func PaginationWithRange(start, end uint64) *PaginationOption {
	return Pagination(start, end-start+1)
}

// CompileSelect to compile select query for pagination
func (p *PaginationOption) CompileSelect(base sq.SelectBuilder) (sq.SelectBuilder, error) {
	if p.limit < 1 {
		return base, errors.New("Limit can't be 0 or negative")
	}
	base = base.Offset(p.offset)
	if p.limit != 0 {
		base = base.Limit(p.limit)
	}
	return base, nil
}
