package dbkit

import (
	sq "github.com/Masterminds/squirrel"
)

// FindOption to compile find query
type FindOption interface {
	CompileQuery(sq.SelectBuilder) (sq.SelectBuilder, error)
}

// FindOptionImpl implementation of FindOption
type FindOptionImpl struct {
	Pagination *Pagination
	sorts      []*Sort
}

// CreateFindOption to create new instance of FindOption
func CreateFindOption() *FindOptionImpl {
	return &FindOptionImpl{}
}

// WithPagination return FindOption with pagination
func (f *FindOptionImpl) WithPagination(pagination *Pagination) *FindOptionImpl {
	f.Pagination = pagination
	return f
}

// WithSort return FindOption with sort
func (f *FindOptionImpl) WithSort(sorts ...*Sort) *FindOptionImpl {
	f.sorts = append(f.sorts, sorts...)
	return f
}

// CompileQuery new select query based on current option
func (f *FindOptionImpl) CompileQuery(base sq.SelectBuilder) (compiled sq.SelectBuilder, err error) {
	compiled = base

	if f.Pagination != nil {
		if compiled, err = f.Pagination.CompileQuery(compiled); err != nil {
			return
		}
	}

	for _, sort := range f.sorts {
		if compiled, err = sort.CompileQuery(compiled); err != nil {
			return
		}
	}
	return
}
