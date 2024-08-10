package charthtml

import (
	"fmt"
	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
	"github.com/tecnologer/code-stats/pkg/models"
	"os"
	"time"
)

const dateFormat = "Jan-02-2006"

func Draw(stats *models.StatsCollection, statType models.StatType, languages ...string) error {
	// create a new bar instance
	bar := charts.NewLine()

	// set some global options like Title/Legend/ToolTip or anything else
	bar.SetGlobalOptions(charts.WithTitleOpts(opts.Title{
		Title:    fmt.Sprintf("%s per Language", statType.Title()),
		Subtitle: getSubtitle(stats, statType),
	}))

	xAxis := make([]string, 0)
	data := make(map[string][]opts.LineData)
	symbols := []string{"circle", "rect", "roundRect", "triangle", "diamond", "pin", "arrow"}
	symbolsForLanguage := make(map[string]string)

	for i, key := range stats.KeysSorted() {
		hasDateDate := false

		for _, stats := range stats.Get(key) {
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
				Name:   stats.Name,
				Value:  stats.ValueOf(statType),
				Symbol: symbol,
			})

			hasDateDate = true
		}

		if hasDateDate {
			xAxis = append(xAxis, key.Format(dateFormat))
		}
	}

	for language, values := range data {
		bar.AddSeries(language, values)
	}

	// Put data into instance
	bar.SetXAxis(xAxis)

	// Where the magic happens
	f, _ := os.Create("bar.html")

	err := bar.Render(f)
	if err != nil {
		return fmt.Errorf("failed to render chart: %w", err)
	}

	return nil
}

func getSubtitle(stats *models.StatsCollection, statType models.StatType) string {
	first := stats.FirstKey()
	last := stats.LastKey()

	if first.Format(time.DateOnly) == last.Format(time.DateOnly) {
		return "at the date " + first.Format(time.DateOnly)
	}

	return fmt.Sprintf("from %s to %s", first.Format(dateFormat), last.Format(dateFormat))
}
