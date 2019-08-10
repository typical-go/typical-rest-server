package typictx

// Release setting
type Release struct {
	Alpha  bool
	GoOS   []string
	GoArch []string
	Github *Github
}

// Github contain github information
type Github struct {
	Owner    string
	RepoName string
}
