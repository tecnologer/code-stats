package cli

import (
	"github.com/urfave/cli/v2"
	"tecnologer.net/code-stats/cmd/flags"
)

func (c *CLI) RunCommand() *cli.Command {
	cmdFlags := flags.Collect()

	cmd := &cli.Command{
		Name:   "run",
		Usage:  "Run the code-stats",
		Action: c.cmdAction,
		Flags:  cmdFlags,
		Subcommands: []*cli.Command{
			c.ExampleCommand(),
		},
	}

	return cmd
}

func (c *CLI) cmdRunAction(ctx *cli.Context) error {
	return nil
}
