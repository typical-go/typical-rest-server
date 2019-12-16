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
	Name string
	Type string
}

// JsonedName is name for json response
func (f Field) JsonedName() string {
	return f.Name
}

// StructTag of Field
func (f Field) StructTag() string {
	return ``
}
