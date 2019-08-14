package typictx

import "github.com/urfave/cli"

// Command represent the command in CLI
type Command struct {
	Name        string
	ShortName   string
	Usage       string
	BeforeFunc  func() error
	ActionFunc  ActionFunc
	SubCommands []*Command
}

// ConvertToCLICommand to convert command ti cli.Command
func ConvertToCLICommand(ctx Context, cmd *Command) (cliCmd cli.Command) {
	cliCmd.Name = cmd.Name
	cliCmd.ShortName = cmd.ShortName
	cliCmd.Usage = cmd.Usage
	if cmd.ActionFunc != nil {
		cliCmd.Action = func(cliCtx *cli.Context) error {
			return cmd.ActionFunc(&ActionContext{
				Context: ctx,
				Cli:     cliCtx,
			})
		}
	}
	if cmd.BeforeFunc != nil {
		cliCmd.Before = func(ctx *cli.Context) error {
			return cmd.BeforeFunc()
		}
	}
	for _, subCmd := range cmd.SubCommands {
		cliCmd.Subcommands = append(cliCmd.Subcommands, ConvertToCLICommand(ctx, subCmd))
	}
	return
}
