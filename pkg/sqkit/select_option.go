package sqkit

import (
	sq "github.com/Masterminds/squirrel"
)

type (
	// SelectOption to compile select query
	SelectOption interface {
		CompileSelect(sq.SelectBuilder) sq.SelectBuilder
	}
	// CompileSelectFn function
	CompileSelectFn  func(sq.SelectBuilder) sq.SelectBuilder
	selectOptionImpl struct {
		fn CompileSelectFn
	}
)

// NewSelectOption return new instance of SelectOption
func NewSelectOption(fn CompileSelectFn) SelectOption {
	return &selectOptionImpl{fn: fn}
}

func (s *selectOptionImpl) CompileSelect(b sq.SelectBuilder) sq.SelectBuilder {
	return s.fn(b)
}
