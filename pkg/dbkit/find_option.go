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
	// Pagination *Pagination
	// sorts      []*Sort
	options []FindOption
}

// CreateFindOption to create new instance of FindOption
func CreateFindOption() *FindOptionImpl {
	return &FindOptionImpl{}
}

func (f *FindOptionImpl) With(option ...FindOption) *FindOptionImpl {
	f.options = append(f.options, option...)
	return f
}

// CompileQuery new select query based on current option
func (f *FindOptionImpl) CompileQuery(base sq.SelectBuilder) (compiled sq.SelectBuilder, err error) {
	compiled = base
	for _, opt := range f.options {
		if compiled, err = opt.CompileQuery(compiled); err != nil {
			return
		}
	}
	return
}
