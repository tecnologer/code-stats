package scc

import (
	"fmt"
	"os"
	"path"

	"github.com/boyter/scc/processor"
	"github.com/tecnologer/code-stats/pkg/file"
	"github.com/tecnologer/code-stats/ui"
)

func Process(omitDir ...string) ([]byte, error) {
	processor.FileOutput = path.Join(os.TempDir(), "updated_stats.json")

	processor.DirFilePaths = []string{}
	if processor.ConfigureLimits != nil {
		processor.ConfigureLimits()
	}

	processor.PathDenyList = omitDir
	processor.Format = "json"
	processor.IgnoreMinified = true

	processor.ConfigureGc()
	processor.ConfigureLazy(true)
	processor.Process()

	data, err := file.ReadContent(processor.FileOutput)
	if err != nil {
		return nil, fmt.Errorf("failed to read stats from tmp file: %w", err)
	}

	err = os.Remove(processor.FileOutput)
	if err != nil {
		ui.Infof("failed to remove tmp file: %s", err)
	}

	ui.Debugf("tmp file removed")

	return data, nil
}
