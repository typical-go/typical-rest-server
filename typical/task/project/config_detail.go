package project

import (
	"bytes"

	"github.com/olekukonko/tablewriter"
	"github.com/typical-go/typical-rest-server/config"
)

// ConfigDetail return config detail string
func ConfigDetail() string {
	buf := new(bytes.Buffer)
	table := tablewriter.NewWriter(buf)
	table.SetHeader([]string{"Name", "Type", "Required", "Default"})
	table.SetBorders(tablewriter.Border{Left: true, Top: false, Right: true, Bottom: false})
	table.SetCenterSeparator("|")
	for _, detail := range config.Informations() {
		table.Append([]string{detail.Name, detail.Type, detail.Required, detail.Default})
	}
	table.Render()

	return buf.String()
}
