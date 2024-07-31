package chart

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/wcharczuk/go-chart/v2"
	"tecnologer.net/code-stats/pkg/models"
	"tecnologer.net/code-stats/ui"
)

func Draw(stats *models.StatsCollection, statType models.StatType, languages ...string) error {
	graph := chart.Chart{
		XAxis: chart.XAxis{
			Name:           "Date",
			ValueFormatter: chart.TimeDateValueFormatter,
		},
		YAxis: chart.YAxis{
			Name: "Code Count",
		},
		Series: []chart.Series{},
	}

	series := map[string]*chart.TimeSeries{}

	for _, key := range stats.Keys() {
		languageStats := stats.Get(key)

		for _, stat := range languageStats {
			if !isInLanguageList(stat.Name, languages) {
				continue
			}

			serie, ok := series[stat.Name]
			if !ok {
				serie = createSeriePerLanguage(languageStats, statType, stat.Name)
				serie.XValues = []time.Time{key}
				series[stat.Name] = serie
			} else {
				serie.XValues = append(serie.XValues, key)
				serie.YValues = append(serie.YValues, float64(stat.ValueOf(statType)))
			}
		}
	}

	for _, serie := range series {
		graph.Series = append(graph.Series, serie)
	}

	graph.Elements = []chart.Renderable{
		chart.Legend(&graph),
	}

	graphFile, err := os.Create(time.Now().UTC().Format(time.DateOnly) + "_stats.png")
	if err != nil {
		return fmt.Errorf("failed to create image file: %w", err)
	}

	defer func(graphFile *os.File) {
		err := graphFile.Close()
		if err != nil {
			ui.Errorf("failed to close image file: %v", err)
		}
	}(graphFile)

	err = graph.Render(chart.PNG, graphFile)
	if err != nil {
		return fmt.Errorf("failed to render chart: %w", err)
	}

	return nil
}

func isInLanguageList(language string, list []string) bool {
	if len(list) == 0 {
		return true
	}

	for _, l := range list {
		if strings.EqualFold(l, language) {
			return true
		}
	}

	return false
}

func createSeriePerLanguage(stats []*models.Stats, statType models.StatType, language string) *chart.TimeSeries {
	serie := &chart.TimeSeries{
		Name:    language,
		YValues: []float64{},
	}

	for _, st := range stats {
		if strings.EqualFold(st.Name, language) {
			serie.YValues = append(serie.YValues, float64(st.ValueOf(statType)))
			break
		}
	}

	return serie
}
