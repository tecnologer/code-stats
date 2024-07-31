package main

import (
	"os"

	"tecnologer.net/code-stats/cmd/cli"
	"tecnologer.net/code-stats/cmd/vers"
	"tecnologer.net/code-stats/ui"
)

var version string

func main() {
	newCLI := cli.NewCLI(vers.Version(version))
	if err := newCLI.Run(os.Args); err != nil {
		ui.Errorf(err.Error())
	}
}
