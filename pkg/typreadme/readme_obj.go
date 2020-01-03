package typreadme

// ReadmeObject represent readme documentation
type ReadmeObject struct {
	Title               string
	Description         string
	ApplicationCommands CommandInfos
	OtherBuildCommands  CommandInfos
	Configs             ConfigInfos
}
