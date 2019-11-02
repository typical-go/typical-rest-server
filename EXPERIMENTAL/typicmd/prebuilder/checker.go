package prebuilder

type checker struct {
	mockTarget      bool
	constructor     bool
	configuration   bool
	testTarget      bool
	buildToolBinary bool
	contextChecksum bool
}

func (r *checker) checkBuildTool() bool {
	return r.mockTarget ||
		r.constructor ||
		r.configuration ||
		r.testTarget ||
		r.buildToolBinary ||
		r.contextChecksum
}
