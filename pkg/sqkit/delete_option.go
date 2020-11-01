package sqkit

import (
	sq "github.com/Masterminds/squirrel"
)

type (
	// DeleteOption to compile delete query
	DeleteOption interface {
		CompileDelete(sq.DeleteBuilder) sq.DeleteBuilder
	}
	// CompileDeleteFn function
	CompileDeleteFn  func(sq.DeleteBuilder) sq.DeleteBuilder
	deleteOptionImpl struct {
		fn CompileDeleteFn
	}
)

// NewDeleteOption return new instance of DeleteOption
func NewDeleteOption(fn CompileDeleteFn) DeleteOption {
	return &deleteOptionImpl{fn: fn}
}

func (s *deleteOptionImpl) CompileDelete(b sq.DeleteBuilder) sq.DeleteBuilder {
	return s.fn(b)
}
