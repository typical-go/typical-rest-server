package golang

// Imports is slice of import
type Imports []Import

// Import is plain old go object for import
type Import struct {
	Name string
	Path string
}

// AddImport to add new import with alias
func (i *Imports) AddImport(name, path string) {
	*i = append(*i, Import{
		Name: name,
		Path: path,
	})
}
