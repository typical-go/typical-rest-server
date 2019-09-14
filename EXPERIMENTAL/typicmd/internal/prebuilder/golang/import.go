package golang

// Imports is slice of import
type Imports []Import

// Import is plain old go object for import
type Import struct {
	Alias       string
	PackageName string
}

// WithAlias to add new import with alias
func (i *Imports) WithAlias(alias, pkg string) {
	*i = append(*i, Import{
		Alias:       alias,
		PackageName: pkg,
	})
}

// Add to add new import
func (i *Imports) Add(pkg string) {
	i.WithAlias("", pkg)
}

// Blank to add new blank import
func (i *Imports) Blank(pkg string) {
	i.WithAlias("_", pkg)
}
