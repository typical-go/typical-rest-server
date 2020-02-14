package typreadme

// Object represent readme documentation
type Object struct {
	TemplateFile  string
	Title         string
	Description   string
	Usages        []CommandInfo
	BuildCommands []CommandInfo
	Configs       []ConfigInfo
}

// CommandInfo is command information
type CommandInfo struct {
	Command     string
	Description string
}

// ConfigInfo is configuration information
type ConfigInfo struct {
	Name     string
	Type     string
	Default  string
	Required string
}
