package typiparser

// ProjectContext is project context
type ProjectContext struct {
	// Layouts contain project layout and folders
	Layouts []string

	// Autowires is list of function what eligible to automatic add to DI
	Autowires []string

	// Automocks is list of file path to be generated the mock
	Automocks []string
}
