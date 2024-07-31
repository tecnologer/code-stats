package flags

import (
	"github.com/urfave/cli/v2"
	"tecnologer.net/code-stats/models"
)

func Verbose() *cli.BoolFlag {
	return &cli.BoolFlag{
		Name:  "verbose",
		Usage: "enable verbose output.",
	}
}

func OmitDirs() *cli.StringSliceFlag {
	return &cli.StringSliceFlag{
		Name:        "omit-dir",
		Usage:       "directories to omit from the stats collection.",
		DefaultText: ".idea,vendor,.stats",
	}
}

/*inputFilePaths   = flag.String("input", ".stats", "Path to the input file, separated by commas, could be used with stdin")
imitDir          = flag.String("omit-dir", ".idea,vendor,.stats", "Directories to omit from the stats")
onlyCompareInput = flag.Bool("only-compare-input", false, "Only compare the input files, do not calculate the current stats")
drawChart        = flag.Bool("draw-chart", false, "Draw chart")
languages        = flag.String("languages", "", "Languages to include in the chart, require at least one and --draw-chart")
statName         = flag.String("stat-name", "code", "Name of the stat, accepted values: "+models.AllStatTypesString())
showVersion      = flag.Bool("version", false, "Show version")
*/

func InputPaths() *cli.StringSliceFlag {
	return &cli.StringSliceFlag{
		Name:        "input",
		Usage:       "list path to the input files or directories",
		DefaultText: ".stats",
	}
}

func OnlyCompareInput() *cli.BoolFlag {
	return &cli.BoolFlag{
		Name:    "only-compare-input",
		Aliases: []string{"oc"},
		Usage:   "only compare the input files, do not calculate the current stats",
	}
}

func DrawChart() *cli.BoolFlag {
	return &cli.BoolFlag{
		Name:    "draw-chart",
		Usage:   "draw chart",
		Aliases: []string{"d"},
	}
}

func Languages() *cli.StringSliceFlag {
	return &cli.StringSliceFlag{
		Name:        "languages",
		Usage:       "languages to include in the chart, require at least one if --draw-chart is set",
		Aliases:     []string{"l"},
		DefaultText: "go",
	}
}

func StatName() *cli.StringFlag {
	return &cli.StringFlag{
		Name:        "stat-name",
		Usage:       "name of the stat, accepted values: " + models.AllStatTypesString(),
		Aliases:     []string{"n"},
		DefaultText: "code",
	}
}

func Collect() []cli.Flag {
	return []cli.Flag{
		OmitDirs(),
	}
}

func Compare() []cli.Flag {
	return []cli.Flag{
		InputPaths(),
		OnlyCompareInput(),
		DrawChart(),
		Languages(),
		StatName(),
	}
}
