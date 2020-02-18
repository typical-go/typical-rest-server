package dbkit

import sq "github.com/Masterminds/squirrel"

// FindOption to compile find query
type FindOption interface {
	CompileQuery(sq.SelectBuilder) (sq.SelectBuilder, error)
	String() string
}
