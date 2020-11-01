package sqkit

import sq "github.com/Masterminds/squirrel"

type (
	// Where conditions as squirel query builder
	Where []interface{}
)

var _ SelectOption = (Eq)(nil)
var _ UpdateOption = (Eq)(nil)
var _ DeleteOption = (Eq)(nil)

// CompileSelect to compile select query for filtering
func (e Where) CompileSelect(base sq.SelectBuilder) sq.SelectBuilder {
	for _, cond := range e {
		base = base.Where(cond)
	}
	return base
}

// CompileUpdate to compile update query for filtering
func (e Where) CompileUpdate(base sq.UpdateBuilder) sq.UpdateBuilder {
	for _, cond := range e {
		base = base.Where(cond)
	}
	return base
}

// CompileDelete to compile delete query for filtering
func (e Where) CompileDelete(base sq.DeleteBuilder) sq.DeleteBuilder {
	for _, cond := range e {
		base = base.Where(cond)
	}
	return base
}
