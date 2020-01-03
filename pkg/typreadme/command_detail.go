package typreadme

// CommandDetails is slice of CommandDetail
type CommandDetails []*CommandDetail

// CommandDetail is detail of command
type CommandDetail struct {
	Snippet string
	Usage   string
}

// Add Command Detail
func (c *CommandDetails) Add(detail *CommandDetail) {
	*c = append(*c, detail)
}
