package sqkit

import (
	sq "github.com/Masterminds/squirrel"
)

type (
	// UpdateOption to compile update query
	UpdateOption interface {
		CompileUpdate(sq.UpdateBuilder) sq.UpdateBuilder
	}
	// CompileUpdateFn function
	CompileUpdateFn  func(sq.UpdateBuilder) sq.UpdateBuilder
	updateOptionImpl struct {
		fn CompileUpdateFn
	}
)

// NewUpdateOption return new instance of UpdateOption
func NewUpdateOption(fn CompileUpdateFn) UpdateOption {
	return &updateOptionImpl{fn: fn}
}

func (u *updateOptionImpl) CompileUpdate(b sq.UpdateBuilder) sq.UpdateBuilder {
	return u.fn(b)
}
