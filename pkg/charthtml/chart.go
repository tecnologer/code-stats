package charthtml

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
	"github.com/tecnologer/code-stats/pkg/models"
	"github.com/tecnologer/code-stats/ui"
)

const dateFormat = "Jan-02-2006"

type DrawOptions struct {
	StatType        models.StatType
	DiffPivot       time.Time
	Languages       []string
	Collection      *models.StatsCollection
	DiffType        models.DifferenceType
	OutputChartPath string
}

func Draw(dOpts *DrawOptions) error {
	// create a new bar instance
	bar := charts.NewLine()

	// set some global options like Title/Legend/ToolTip or anything else
	bar.SetGlobalOptions(charts.WithTitleOpts(opts.Title{
		Title:    getTitle(dOpts),
		Subtitle: getSubtitle(dOpts),
	}))

	xAxis := make([]string, 0)
	data := make(map[string][]opts.LineData)
	symbols := NewSymbol()

	for _, key := range dOpts.Collection.KeysSorted() {
		hasDateDate := false

		for _, stats := range dOpts.Collection.Get(key) {
			if !stats.IsInLanguageList(dOpts.Languages) {
				continue
			}

			if _, ok := data[stats.Name]; !ok {
				data[stats.Name] = make([]opts.LineData, 0, 1)
			}

			data[stats.Name] = append(data[stats.Name], opts.LineData{
				Name:       stats.Name,
				Value:      getValueFor(dOpts, key, stats),
				Symbol:     symbols.GetFor(stats.Name),
				SymbolSize: 10,
			})

			hasDateDate = true
		}

		if hasDateDate {
			xAxis = append(xAxis, key.Format(dateFormat))
		}
	}

	for language, values := range data {
		bar.AddSeries(language, values, charts.WithLabelOpts(opts.Label{
			Show: opts.Bool(true),
		}))
	}

	// Put data into instance
	bar.SetXAxis(xAxis)

	err := drawFile(bar, dOpts.OutputChartPath)
	if err != nil {
		return fmt.Errorf("failed to save chart to file: %w", err)
	}

	return nil
}

func getTitle(dOpts *DrawOptions) string {
	if dOpts.DiffType != models.DiffNone {
		return fmt.Sprintf("Progress of %s per Language", dOpts.StatType.Title())
	}

	return dOpts.StatType.Title() + " per Language"
}

func getSubtitle(dOpts *DrawOptions) string {
	first := dOpts.Collection.FirstKey()
	last := dOpts.Collection.LastKey()

	var subtitle string

	if dOpts.DiffType != models.DiffNone {
		const differenceSubtitleFmt = "difference calculated against %s, "

		switch {
		case dOpts.DiffType == models.DiffPreviousDate:
			subtitle = fmt.Sprintf(differenceSubtitleFmt, dOpts.DiffType)
		case dOpts.DiffType == models.DiffFirstDate:
			subtitle = fmt.Sprintf(differenceSubtitleFmt, first.Format(dateFormat)+" (first date)")
		default:
			subtitle = fmt.Sprintf(differenceSubtitleFmt, dOpts.DiffPivot.Format(dateFormat))
		}
	}

	if first.Format(time.DateOnly) == last.Format(time.DateOnly) {
		subtitle += "at the date " + first.Format(time.DateOnly)
	} else {
		subtitle += fmt.Sprintf("from %s to %s", first.Format(dateFormat), last.Format(dateFormat))
	}

	return subtitle
}

func drawFile(graph *charts.Line, outputChartPath string) error {
	if outputChartPath == "" {
		outputChartPath = time.Now().UTC().Format(time.DateOnly) + "_stats.html"
	}

	if !strings.HasSuffix(outputChartPath, ".html") {
		outputChartPath += ".html"
	}

	graphFile, err := os.Create(outputChartPath)
	if err != nil {
		return fmt.Errorf("failed to create image file: %w", err)
	}

	defer func(graphFile *os.File) {
		err := graphFile.Close()
		if err != nil {
			ui.Errorf("failed to close image file: %v", err)
		}
	}(graphFile)

	err = graph.Render(graphFile)
	if err != nil {
		return fmt.Errorf("failed to render chart: %w", err)
	}

	return nil
}

func getValueFor(dOpts *DrawOptions, currentKey time.Time, stats *models.Stats) float64 {
	if dOpts.DiffType != models.DiffNone {
		return dOpts.Collection.DiffPrevious(currentKey, stats.Name, dOpts.StatType, dOpts.DiffPivot)
	}

	return stats.ValueOf(dOpts.StatType)
}
