package typrest

// Entity of database
type Entity struct {
	Name           string
	Type           string
	Table          string
	Cache          string
	ProjectPackage string
	Fields         []Field
}

// Field of entity
type Field struct {
	Name   string
	Type   string
	Column string
}

// StructTag of Field
func (f Field) StructTag() string {
	return ``
}
