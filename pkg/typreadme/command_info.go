package typreadme

// CommandInfos is collection of command  information
type CommandInfos []*CommandInfo

// CommandInfo is command information
type CommandInfo struct {
	Snippet string
	Usage   string
}

// Append command info
func (c *CommandInfos) Append(detail *CommandInfo) {
	*c = append(*c, detail)
}
