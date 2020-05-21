package dbkit

import sq "github.com/Masterminds/squirrel"

type (
	// SelectOption to compile find query
	SelectOption interface {
		CompileSelect(sq.SelectBuilder) (sq.SelectBuilder, error)
	}
)
