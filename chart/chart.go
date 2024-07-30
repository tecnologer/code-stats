package chart

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/wcharczuk/go-chart/v2"
	"tecnologer.net/code-stats/models"
)

func Draw(stats *models.StatsCollection, languages ...string) error {
	//bars := make([]chart.Value, 0, stats.Len())

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
				serie = createSeriePerLanguage(languageStats, stat.Name)
				serie.XValues = stats.Keys()
				series[stat.Name] = serie
			} else {
				serie.YValues = append(serie.YValues, float64(stat.Code))
			}
		}
	}

	for _, serie := range series {
		graph.Series = append(graph.Series, serie)
	}

	//graph := chart.BarChart{
	//	Title: "Code stats",
	//	Background: chart.Style{
	//		Padding: chart.Box{
	//			Top: 40,
	//		},
	//	},
	//	Height:   512,
	//	Width:    512,
	//	BarWidth: 60,
	//	Bars:     bars,
	//}

	f, _ := os.Create(time.Now().UTC().Format(time.DateOnly) + "_stats.png")
	defer f.Close()

	err := graph.Render(chart.PNG, f)
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

func createSeriePerLanguage(stats []*models.Stats, language string) *chart.TimeSeries {
	serie := &chart.TimeSeries{
		Name:    language,
		YValues: []float64{},
	}

	for _, st := range stats {
		if strings.EqualFold(st.Name, language) {
			serie.YValues = append(serie.YValues, float64(st.Code))
			break
		}
	}

	return serie
}
