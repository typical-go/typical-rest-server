package typiast

// Report of AST transversal
type Report struct {
	Packages []string `json:"packages"`
	Files    []File   `json:"files"`
}

// AddFile to add new file
func (r *Report) AddFile(file File) *Report {
	if !file.IsEmpty() {
		r.Files = append(r.Files, file)
	}
	return r
}

// Autowires return autowired constructors
func (r *Report) Autowires() (constructors []string) {
	for _, file := range r.Files {
		constructors = append(constructors, file.Constructors...)
	}
	return
}

// Automocks return automocked filenames
func (r *Report) Automocks() (filenames []string) {
	for _, file := range r.Files {
		if file.Mock {
			filenames = append(filenames, file.Name)
		}
	}
	return
}
