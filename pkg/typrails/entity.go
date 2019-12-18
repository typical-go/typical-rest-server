package typrails

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
