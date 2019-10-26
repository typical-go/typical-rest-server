package typictx

import "github.com/urfave/cli"

// CommandLiner responsible to give command
type CommandLiner interface {
	CommandLine() *Command
}

// Command represent the command in CLI
type Command struct {
	Name        string
	ShortName   string
	Usage       string
	BeforeFunc  func() error
	Flags       []cli.Flag
	ActionFunc  ActionFunc
	SubCommands []*Command
}

// CliCommand to convert command ti cli.Command
func (c *Command) CliCommand(ctx *Context) (cliCmd cli.Command) {
	cliCmd.Name = c.Name
	cliCmd.ShortName = c.ShortName
	cliCmd.Usage = c.Usage
	cliCmd.Flags = c.Flags
	if c.ActionFunc != nil {
		cliCmd.Action = func(cliCtx *cli.Context) error {
			return c.ActionFunc(&ActionContext{
				Context: ctx,
				Cli:     cliCtx,
			})
		}
	}
	if c.BeforeFunc != nil {
		cliCmd.Before = func(ctx *cli.Context) error {
			return c.BeforeFunc()
		}
	}
	for _, subCmd := range c.SubCommands {
		cliCmd.Subcommands = append(cliCmd.Subcommands, subCmd.CliCommand(ctx))
	}
	return
}
