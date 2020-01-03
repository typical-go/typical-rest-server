package typreadme

// ReadmeObject represent readme documentation
type ReadmeObject struct {
	Template            string
	Title               string
	Description         string
	ApplicationCommands CommandInfos
	OtherBuildCommands  CommandInfos
	Configs             ConfigInfos
}
