package walker

// ProjectFiles information
type ProjectFiles []ProjectFile

// Add item to files
func (f *ProjectFiles) Add(item ProjectFile) {
	*f = append(*f, item)
}

// Autowires return autowired constructors
func (f *ProjectFiles) Autowires() (constructors []string) {
	for _, file := range *f {
		constructors = append(constructors, file.Constructors...)
	}
	return
}

// Automocks return automocked filenames
func (f *ProjectFiles) Automocks() (filenames []string) {
	for _, file := range *f {
		if file.Mock {
			filenames = append(filenames, file.Name)
		}
	}
	return
}
