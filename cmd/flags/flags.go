package flags

import (
	"github.com/urfave/cli/v2"
	"tecnologer.net/code-stats/pkg/file"
	"tecnologer.net/code-stats/pkg/models"
)

const (
	VerboseFlagName          = "verbose"
	OmitDirsFlagName         = "omit-dir"
	InputPathsFlagName       = "input"
	OnlyCompareInputFlagName = "only-compare-input"
	DrawChartFlagName        = "draw-chart"
	LanguagesFlagName        = "languages"
	StatNameFlagName         = "stat-name"
	NoEmojisFlagName         = "no-emoji"
	NoColorFlagName          = "no-color"
)

func Verbose() *cli.BoolFlag {
	return &cli.BoolFlag{
		Name:  VerboseFlagName,
		Usage: "enable verbose output.",
	}
}

func NoColor() *cli.BoolFlag {
	return &cli.BoolFlag{
		Name:  NoColorFlagName,
		Usage: "disable color output.",
	}
}

func OmitDirs() *cli.StringSliceFlag {
	return &cli.StringSliceFlag{
		Name:  OmitDirsFlagName,
		Usage: "directories to omit from the stats collection.",
		Value: cli.NewStringSlice(".idea", "vendor", ".stats"),
	}
}

func InputPaths() *cli.StringSliceFlag {
	return &cli.StringSliceFlag{
		Name:  InputPathsFlagName,
		Usage: "list path to the input files or directories",
		Value: cli.NewStringSlice(file.StatsDirectoryPath),
	}
}

func OnlyCompareInput() *cli.BoolFlag {
	return &cli.BoolFlag{
		Name:    OnlyCompareInputFlagName,
		Aliases: []string{"c"},
		Usage:   "only compare the input files, do not calculate the current stats",
	}
}

func DrawChart() *cli.BoolFlag {
	return &cli.BoolFlag{
		Name:    DrawChartFlagName,
		Usage:   "draw chart",
		Aliases: []string{"d"},
	}
}

func Languages() *cli.StringSliceFlag {
	return &cli.StringSliceFlag{
		Name:    LanguagesFlagName,
		Usage:   "languages to include in the chart, require at least one if --draw-chart is set",
		Aliases: []string{"l"},
		Value:   cli.NewStringSlice("go"),
	}
}

func StatName() *cli.StringFlag {
	return &cli.StringFlag{
		Name:    StatNameFlagName,
		Usage:   "name of the stat, accepted values: " + models.AllStatTypesString(),
		Aliases: []string{"n"},
		Value:   "code",
	}
}

func NoEmojis() *cli.BoolFlag {
	return &cli.BoolFlag{
		Name:  NoEmojisFlagName,
		Usage: "disable emojis in the output.",
	}
}
