package dbkit

import sq "github.com/Masterminds/squirrel"

type (
	// SelectOption to compile select query
	SelectOption interface {
		CompileSelect(sq.SelectBuilder) (sq.SelectBuilder, error)
	}

	// UpdateOption to compile update query
	UpdateOption interface {
		CompileUpdate(sq.UpdateBuilder) (sq.UpdateBuilder, error)
	}

	// DeleteOption to compile delete query
	DeleteOption interface {
		CompileDelete(sq.DeleteBuilder) (sq.DeleteBuilder, error)
	}
)
