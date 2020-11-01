package sqkit

import (
	"fmt"
	"strings"

	sq "github.com/Masterminds/squirrel"
)

type (
	// Sorts sorting
	Sorts []string
)

//
// Sort
//

var _ SelectOption = (*Sorts)(nil)

// CompileSelect to compile select query for sorting
func (s Sorts) CompileSelect(base sq.SelectBuilder) sq.SelectBuilder {
	for _, sort := range s {
		base = base.OrderBy(s.statement(sort))
	}
	return base
}

func (s Sorts) statement(sort string) string {
	var column, orderBy string
	if strings.HasPrefix(sort, "-") {
		column = sort[1:]
		orderBy = "DESC"
	} else if strings.HasPrefix(sort, "+") {
		column = sort[1:]
		orderBy = "ASC"
	} else {
		column = sort
		orderBy = "ASC"
	}
	return fmt.Sprintf("%s %s", column, orderBy)
}
