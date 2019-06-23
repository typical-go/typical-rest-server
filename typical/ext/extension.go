package ext

import "gopkg.in/urfave/cli.v1"

// Extension representated typical setup and command
type Extension interface {
	Setup() error
	Command() cli.Command
}
