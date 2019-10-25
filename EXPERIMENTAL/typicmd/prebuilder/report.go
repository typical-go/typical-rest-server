package prebuilder

type report struct {
	MockTargetUpdated    bool
	ConstructorUpdated   bool
	ConfigurationUpdated bool
	TestTargetUpdated    bool
}

func (r *report) Updated() bool {
	return r.MockTargetUpdated ||
		r.ConstructorUpdated ||
		r.ConfigurationUpdated ||
		r.TestTargetUpdated
}
