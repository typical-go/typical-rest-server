package walker

// ProjectFile of walk analysis
type ProjectFile struct {
	Name         string
	Mock         bool
	Constructors []string
}

// IsEmpty is true if empty truct
func (f *ProjectFile) IsEmpty() bool {
	return !f.Mock && len(f.Constructors) < 1
}

// AddConstructor to add contructor to file
func (f *ProjectFile) AddConstructor(constructor string) *ProjectFile {
	f.Constructors = append(f.Constructors, constructor)
	return f
}
