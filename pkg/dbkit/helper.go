package dbkit

import sq "github.com/Masterminds/squirrel"

// CompileOpts the find options
func CompileOpts(builder sq.SelectBuilder, opts ...FindOption) (err error) {
	for _, opt := range opts {
		if builder, err = opt.CompileQuery(builder); err != nil {
			return
		}
	}
	return
}
