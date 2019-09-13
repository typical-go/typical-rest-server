package walker

// File of walk analysis
type File struct {
	Name         string
	Mock         bool
	Constructors []string
}

// IsEmpty is true if empty truct
func (f *File) IsEmpty() bool {
	return !f.Mock && len(f.Constructors) < 1
}

// AddConstructor to add contructor to file
func (f *File) AddConstructor(constructor string) *File {
	f.Constructors = append(f.Constructors, constructor)
	return f
}
