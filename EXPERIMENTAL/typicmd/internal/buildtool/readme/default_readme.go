package readme

// DefaultReadme create default readme
func DefaultReadme() *Readme {
	return NewReadme().
		SetSection("Getting Started", defaultGettingStarted).
		SetSection("Release Distribution", defaultReleaseDistribution)
}

func defaultGettingStarted(md *Markdown) (err error) {
	md.Heading3("Prerequisite")
	md.OrderedList(
		"Install [Go](https://golang.org/doc/install) or `brew install go`",
	)

	md.Heading3("Run")
	md.Writeln("Use `./typicalw run` to compile and run local development. You can find the binary at `bin` folder")
	return
}

func defaultReleaseDistribution(md *Markdown) (err error) {
	md.Writeln("Use `./typicalw release` to make the release. You can find the binary at `release` folder. More information check [here](https://typical-go.github.io/release.html)")
	return
}
