package main

import (
	"os"

	"github.com/tecnologer/code-stats/cmd/cli"
	"github.com/tecnologer/code-stats/cmd/vers"
	"github.com/tecnologer/code-stats/ui"
)

var version string

func main() {
	newCLI := cli.NewCLI(vers.Version(version))
	if err := newCLI.Run(os.Args); err != nil {
		ui.Errorf(err.Error())
	}
}
