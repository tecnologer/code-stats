package cli

import (
	"github.com/urfave/cli/v2"
	"tecnologer.net/code-stats/cmd/flags"
)

type CLI struct {
	*cli.App
}

func NewCLI(versionValue string) *CLI {
	newCLI := &CLI{}

	newCLI.setupApp(versionValue)

	return newCLI
}

func (c *CLI) setupApp(versionValue string) {
	c.App = &cli.App{
		Name:        "code-stats",
		Version:     versionValue,
		Usage:       "Collects the code statistics of a given directory, and could compare with the previous stats.",
		Description: "Code Stats is a tool that collects the code statistics of a given directory, and could compare with the previous stats.",
		Commands: []*cli.Command{
			c.CollectStatsCommand(),
			c.CompareStatsCommand(),
		},
		Flags: []cli.Flag{
			flags.Verbose(),
		},
		EnableBashCompletion: true,
	}
}
