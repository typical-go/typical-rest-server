package sqkit

import (
	sq "github.com/Masterminds/squirrel"
)

type (
	// Eq equal
	Eq map[string]interface{}
)

var _ SelectOption = (Eq)(nil)
var _ UpdateOption = (Eq)(nil)
var _ DeleteOption = (Eq)(nil)

// CompileSelect to compile select query for filtering
func (e Eq) CompileSelect(base sq.SelectBuilder) sq.SelectBuilder {
	if len(e) > 0 {
		return base.Where(sq.Eq(e))
	}
	return base
}

// CompileUpdate to compile update query for filtering
func (e Eq) CompileUpdate(base sq.UpdateBuilder) sq.UpdateBuilder {
	if len(e) > 0 {
		return base.Where(sq.Eq(e))
	}
	return base
}

// CompileDelete to compile delete query for filtering
func (e Eq) CompileDelete(base sq.DeleteBuilder) sq.DeleteBuilder {
	if len(e) > 0 {
		return base.Where(sq.Eq(e))
	}
	return base
}
