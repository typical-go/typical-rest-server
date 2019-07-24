package typitask

import (
	"os"

	"github.com/olekukonko/tablewriter"
	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/typictx"
)

// CheckStatus to check if application ready to start
func CheckStatus(ctx typictx.ActionContext) error {
	statusReport := ctx.Typical.CheckModuleStatus()
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Module Name", "Status"})

	for moduleName, status := range statusReport {
		table.Append([]string{moduleName, status})
	}
	table.Render()
	return nil
}
