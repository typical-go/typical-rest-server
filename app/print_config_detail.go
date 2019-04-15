package app

import (
	"os"

	"github.com/imantung/typical-go-server/config"
	"github.com/olekukonko/tablewriter"
	"github.com/urfave/cli"
)

func printConfigDetails(ctx *cli.Context) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Name", "Type", "Required", "Default"})
	for _, detail := range config.Details() {
		table.Append([]string{detail.Name, detail.Type, detail.Required, detail.Default})
	}
	table.Render()
}
