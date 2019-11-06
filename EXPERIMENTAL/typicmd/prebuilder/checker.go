package prebuilder

type checker struct {
	mockTarget bool
	// constructor     bool
	configuration   bool
	testTarget      bool
	buildToolBinary bool
	contextChecksum bool
	buildCommands   bool
	readmeFile      bool
}

func (r *checker) checkBuildTool() bool {
	return r.mockTarget ||
		r.configuration ||
		r.testTarget ||
		r.buildToolBinary ||
		r.contextChecksum ||
		r.buildCommands
}

func (r *checker) checkReadme() bool {
	return r.buildCommands ||
		r.configuration ||
		r.readmeFile
}
