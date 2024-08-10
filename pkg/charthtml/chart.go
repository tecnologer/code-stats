package charthtml

import (
	"fmt"
	"os"
	"time"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
	"github.com/tecnologer/code-stats/pkg/models"
	"github.com/tecnologer/code-stats/ui"
)

const dateFormat = "Jan-02-2006"

func Draw(collection *models.StatsCollection, statType models.StatType, isDiff bool, languages ...string) error {
	// create a new bar instance
	bar := charts.NewLine()

	// set some global options like Title/Legend/ToolTip or anything else
	bar.SetGlobalOptions(charts.WithTitleOpts(opts.Title{
		Title:    getTitle(statType, isDiff),
		Subtitle: getSubtitle(collection, statType),
	}))

	xAxis := make([]string, 0)
	data := make(map[string][]opts.LineData)
	symbols := []string{"circle", "triangle", "diamond", "rect", "roundRect", "pin", "arrow"}
	symbolsForLanguage := make(map[string]string)

	for i, key := range collection.KeysSorted() {
		hasDateDate := false

		for _, stats := range collection.Get(key) {
			if !stats.IsInLanguageList(languages) {
				continue
			}

			symbol, ok := symbolsForLanguage[stats.Name]
			if !ok {
				symbolsForLanguage[stats.Name] = symbols[i]
			}

			if _, ok := data[stats.Name]; !ok {
				data[stats.Name] = make([]opts.LineData, 0)
			}

			data[stats.Name] = append(data[stats.Name], opts.LineData{
				Name:       stats.Name,
				Value:      getValueFor(collection, key, stats, statType, isDiff),
				Symbol:     symbol,
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

	err := drawFile(bar)
	if err != nil {
		return fmt.Errorf("failed to save chart to file: %w", err)
	}

	return nil
}

func getTitle(statType models.StatType, isDiff bool) string {
	if isDiff {
		return fmt.Sprintf("Progress of %s per Language", statType.Title())
	}

	return fmt.Sprintf("%s per Language", statType.Title())

}

func getSubtitle(stats *models.StatsCollection, statType models.StatType) string {
	first := stats.FirstKey()
	last := stats.LastKey()

	if first.Format(time.DateOnly) == last.Format(time.DateOnly) {
		return "at the date " + first.Format(time.DateOnly)
	}

	return fmt.Sprintf("from %s to %s", first.Format(dateFormat), last.Format(dateFormat))
}

func drawFile(graph *charts.Line) error {
	graphFile, err := os.Create(time.Now().UTC().Format(time.DateOnly) + "_stats.html")
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

func getValueFor(collection *models.StatsCollection, currentKey time.Time, stats *models.Stats, statType models.StatType, isDiff bool) float64 {
	if isDiff {
		return collection.DiffPrevious(currentKey, stats.Name, statType)
	}

	return stats.ValueOf(statType)
}
