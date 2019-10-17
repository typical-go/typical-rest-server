package prebuilder

type report struct {
	AnnotatedUpdated     bool
	ConfigurationUpdated bool
	TestTargetUpdated    bool
}

func (r *report) Updated() bool {
	return r.AnnotatedUpdated ||
		r.ConfigurationUpdated ||
		r.TestTargetUpdated
}
