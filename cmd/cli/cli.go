package cli

import (
	"fmt"
	"github.com/tecnologer/code-stats/pkg/charthtml"

	"github.com/fatih/color"
	"github.com/tecnologer/code-stats/cmd/flags"
	"github.com/tecnologer/code-stats/pkg/extractor"
	"github.com/tecnologer/code-stats/pkg/file"
	"github.com/tecnologer/code-stats/pkg/models"
	"github.com/tecnologer/code-stats/ui"
	"github.com/urfave/cli/v2"
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
		Name:    "code-stats",
		Version: versionValue,
		Usage:   "Collects the code statistics of a given directory, and could draw a chart to compare with the previous stats.",
		Description: "Code Stats is a tool that collects the code statistics of a given directory, and could compare with the previous stats. \n" +
			"Understand as  code statistics the number of lines, files, complexity, and other metrics of the codebase. \n",
		Action: c.run,
		Before: c.beforeRun,
		Flags: []cli.Flag{
			flags.Verbose(),
			flags.NoEmojis(),
			flags.NoColor(),
			flags.OmitDirs(),
			flags.InputPaths(),
			flags.OnlyCompareInput(),
			flags.DrawChart(),
			flags.Languages(),
			flags.StatName(),
		},
		EnableBashCompletion: true,
	}
}

func (c *CLI) beforeRun(ctx *cli.Context) error {
	// Disable color globally.
	if ctx.Bool(flags.NoColorFlagName) {
		color.NoColor = true
	}

	if ctx.Bool(flags.VerboseFlagName) {
		ui.SetOutputLevel(ui.DebugLevel)
	}

	if ctx.Bool(flags.NoEmojisFlagName) {
		ui.SetEmojiVisibility(false)
	}

	return nil
}

func (c *CLI) run(ctx *cli.Context) error {
	stats, err := c.extractData(ctx)
	if err != nil {
		return fmt.Errorf("failed to extract data: %w", err)
	}

	if ctx.Bool(flags.DrawChartFlagName) {
		err := c.drawChart(stats, ctx)
		if err != nil {
			return fmt.Errorf("failed to draw chart: %w", err)
		}
	}

	return nil
}

func (c *CLI) extractData(ctx *cli.Context) (*models.StatsCollection, error) {
	stats := models.NewCollection()

	if ctx.Bool(flags.DrawChartFlagName) {
		err := c.couldDrawChart(ctx)
		if err != nil {
			return nil, fmt.Errorf("the chart could not be drawn: %w", err)
		}

		inputStats, err := extractor.ExtractFromInput(ctx.StringSlice(flags.InputPathsFlagName))
		if err != nil {
			return nil, fmt.Errorf("failed to extract data from inputs: %w", err)
		}

		stats.Merge(inputStats)
	}

	if !ctx.Bool(flags.OnlyCompareInputFlagName) {
		currentStats, err := extractor.ExtractCurrent(ctx.StringSlice(flags.OmitDirsFlagName))
		if err != nil {
			return nil, fmt.Errorf("failed to extract current stats: %w", err)
		}

		stats.Merge(currentStats)
	}

	ui.Infof("stats collected successfully")

	return stats, nil
}

func (c *CLI) drawChart(stats *models.StatsCollection, ctx *cli.Context) error {
	statName := ctx.String(flags.StatNameFlagName)

	statType := models.StatTypeFromString(statName)
	if !statType.IsValid() {
		return fmt.Errorf("invalid stat name: %s. The valid stat names are: %s", statName, models.AllStatTypesString())
	}

	err := charthtml.Draw(stats, statType, ctx.StringSlice(flags.LanguagesFlagName)...)
	if err != nil {
		return fmt.Errorf("failed to draw chart: %w", err)
	}

	ui.Successf("chart generated successfully")

	return nil
}

func (c *CLI) couldDrawChart(ctx *cli.Context) error {
	inputPath := ctx.StringSlice(flags.InputPathsFlagName)

	if len(inputPath) == 0 {
		return fmt.Errorf("no input paths provided")
	}

	for _, p := range inputPath {
		if p == "" {
			return fmt.Errorf("empty input path provided")
		}

		if p == file.StatsDirectoryPath && !file.IsPathExists(p) && len(inputPath) == 1 {
			return fmt.Errorf("there is no data to compare")
		}
	}

	return nil
}
