package walker

// Files information
type Files []File

// Add item to files
func (f *Files) Add(item File) {
	*f = append(*f, item)
}

// Autowires return autowired constructors
func (f *Files) Autowires() (constructors []string) {
	for _, file := range *f {
		constructors = append(constructors, file.Constructors...)
	}
	return
}

// Automocks return automocked filenames
func (f *Files) Automocks() (filenames []string) {
	for _, file := range *f {
		if file.Mock {
			filenames = append(filenames, file.Name)
		}
	}
	return
}
