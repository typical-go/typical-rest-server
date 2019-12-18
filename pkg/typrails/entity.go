package typrails

import (
	"fmt"

	"github.com/iancoleman/strcase"
)

// Entity of database
type Entity struct {
	Name           string
	Type           string
	Table          string
	Cache          string
	ProjectPackage string
	Fields         []Field
	Forms          []Field
}

// Field of entity
type Field struct {
	Name      string
	Type      string
	Udt       string
	Column    string
	StructTag string
}

// CreateField to create field
func CreateField(column, udt string) Field {
	return Field{
		Name:      fieldName(column),
		Type:      fieldType(udt),
		Udt:       udt,
		Column:    column,
		StructTag: fmt.Sprintf("`json:\"%s\"`", column),
	}
}

func fieldName(column string) string {
	if column == "id" {
		return "ID"
	}
	return strcase.ToCamel(column)
}

func fieldType(udt string) string {
	switch udt {
	case "int4":
		return "int64"
	case "varchar":
		return "string"
	case "timestamp":
		return "time.Time"
	}
	return "interface{}"
}
