package cli

import (
	"fmt"
	"time"

	"github.com/fatih/color"
	"github.com/tecnologer/code-stats/cmd/flags"
	"github.com/tecnologer/code-stats/pkg/charthtml"
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
		Description: "Code Stats is a tool that leverages SCC to gather code statistics from a specified directory, " +
			"allowing you to compare these stats over time and visualize your progress with charts.\n" +
			"\n* Use the `--diff` flag to compare the current stats with those from the previous day, the first day, or a specific date." +
			"\n* The `--diff-pivot` flag allows you to specify the date to compare the stats with. By default, this is set to the " +
			"'previous-day' (the last day with recorded data, not necessarily yesterday)." +
			"\n   > You can also use the value `first-date` to compare with the first recorded date in the dataset." +
			"\n   > Dates should be provided in the format 'YYYY-MM-DD'.",
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
			flags.Diff(),
			flags.DiffPivot(),
			flags.OutputChartPath(),
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

		ui.Successf("current stats calculated successfully")

		stats.Merge(currentStats)
	}

	return stats, nil
}

func (c *CLI) drawChart(stats *models.StatsCollection, ctx *cli.Context) error {
	dOpts, err := c.drawOptions(stats, ctx)
	if err != nil {
		return fmt.Errorf("failed to get draw options: %w", err)
	}

	err = charthtml.Draw(dOpts)
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

func (c *CLI) drawOptions(stats *models.StatsCollection, ctx *cli.Context) (*charthtml.DrawOptions, error) {
	statName := ctx.String(flags.StatNameFlagName)

	statType := models.StatTypeFromString(statName)
	if !statType.IsValid() {
		return nil, fmt.Errorf("invalid stat name: %s. The valid stat names are: %s", statName, models.AllStatTypesString())
	}

	calculateDiff := ctx.Bool(flags.CalculateDiffFlagName)
	if calculateDiff {
		ui.Infof("drawing chart with diff")
	}

	diffPivot := ctx.String(flags.CalculateDiffPivotFlagName)
	diffType := c.getDiffType(calculateDiff, diffPivot)

	pivotDiff, err := c.getDiffPivot(diffType, diffPivot, stats)
	if err != nil {
		return nil, fmt.Errorf("failed to get diff pivot: %w", err)
	}

	return &charthtml.DrawOptions{
		StatType:        statType,
		DiffType:        diffType,
		DiffPivot:       pivotDiff,
		Languages:       ctx.StringSlice(flags.LanguagesFlagName),
		OutputChartPath: ctx.String(flags.OutputChartPathFlagName),
		Collection:      stats,
	}, nil
}

func (c *CLI) getDiffType(calculateDiff bool, diffPivot string) models.DifferenceType {
	if !calculateDiff {
		return models.DiffNone
	}

	diff, _ := models.DifferenceTypeString(diffPivot)

	if diffPivot == "" || diff == models.DiffPreviousDate {
		return models.DiffPreviousDate
	}

	if diff == models.DiffFirstDate {
		return models.DiffFirstDate
	}

	return models.DiffSpecificDate
}

func (c *CLI) getDiffPivot(diffType models.DifferenceType, pivot string, collection *models.StatsCollection) (time.Time, error) {
	if diffType == models.DiffNone {
		return time.Time{}, nil
	}

	if diffType == models.DiffPreviousDate {
		ui.Infof("using the previous day as pivot to calculate difference")

		return time.Time{}, nil
	}

	if diffType == models.DiffFirstDate {
		firstDate := collection.FirstKey()

		ui.Infof("using the first date %s as pivot to calculate difference", firstDate.Format(time.DateOnly))

		return firstDate, nil
	}

	pivotTime, err := time.Parse(time.DateOnly, pivot)
	if err != nil {
		return time.Time{}, fmt.Errorf("failed to parse pivot for difference: %w", err)
	}

	if collection.Get(pivotTime) == nil {
		return time.Time{}, fmt.Errorf("the pivot date '%s' does not exist in the stats collection", pivotTime.Format(time.DateOnly))
	}

	ui.Infof("using the date %s as pivot to calculate difference", pivotTime.Format(time.DateOnly))

	return pivotTime, nil
}
