package chart

import (
	"fmt"
	"github.com/wcharczuk/go-chart/v2/drawing"
	"os"
	"strings"
	"time"

	"github.com/tecnologer/code-stats/pkg/models"
	"github.com/tecnologer/code-stats/ui"
	"github.com/wcharczuk/go-chart/v2"
)

func Draw(stats *models.StatsCollection, statType models.StatType, languages ...string) error {
	graph := chart.Chart{
		XAxis: chart.XAxis{
			Name:           "Date",
			ValueFormatter: chart.TimeDateValueFormatter,
		},
		YAxis: chart.YAxis{
			Name: fmt.Sprintf("%s Count", statType),
		},
		Series: []chart.Series{},
	}

	series := map[string]*AnnotatedTimeSeries{}

	for _, key := range stats.Keys() {
		languageStats := stats.Get(key)

		for _, stat := range languageStats {
			if !isInLanguageList(stat.Name, languages) {
				continue
			}

			serie, ok := series[stat.Name]
			if ok {
				serie.XValues = append(serie.XValues, key)
				serie.YValues = append(serie.YValues, float64(stat.ValueOf(statType)))

				continue
			}

			serie = createSeriePerLanguage(languageStats, statType, stat.Name)
			serie.XValues = []time.Time{key}
			series[stat.Name] = serie
		}
	}

	for _, serie := range series {
		graph.Series = append(graph.Series, serie)
	}

	graph.Elements = []chart.Renderable{
		LegendLeft(&graph),
	}

	err := drawFile(&graph)
	if err != nil {
		return fmt.Errorf("failed to save chart to file: %w", err)
	}

	return nil
}

func drawFile(graph *chart.Chart) error {
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

func createSeriePerLanguage(stats []*models.Stats, statType models.StatType, language string) *AnnotatedTimeSeries {
	serie := &AnnotatedTimeSeries{
		TimeSeries: chart.TimeSeries{
			Name:    language,
			YValues: []float64{},
			Style: chart.Style{
				DotWidth: 2.8,
			},
		},
	}

	for _, st := range stats {
		if strings.EqualFold(st.Name, language) {
			serie.YValues = append(serie.YValues, float64(st.ValueOf(statType)))
			break
		}
	}

	return serie
}

type AnnotatedTimeSeries struct {
	chart.TimeSeries
}

// Render renders the series to the chart.
func (ats AnnotatedTimeSeries) Render(r chart.Renderer, canvasBox chart.Box, xrange, yrange chart.Range, defaults chart.Style) {
	// First, render the line as usual
	ats.TimeSeries.Render(r, canvasBox, xrange, yrange, defaults)

	// Now, render annotations for each point
	style := ats.Style.InheritFrom(defaults)
	r.SetFont(style.GetFont())
	r.SetFontColor(style.GetFontColor())

	for index, x := range ats.XValues {
		// Convert the x and y values to their corresponding points on the canvas
		xvValue := xrange.Translate(float64(x.Unix()))
		yValue := yrange.Translate(ats.YValues[index])

		// Format the label text
		label := fmt.Sprintf("%.2f", ats.YValues[index])

		// Calculate text width and height (simple approximation)
		// Calculate text width and height using MeasureText
		textDimensions := r.MeasureText(label)

		// Draw the text slightly above and to the right of the point
		textX := canvasBox.Left + xvValue - textDimensions.Width()/2
		textY := canvasBox.Top + yValue - textDimensions.Height() - 5 // 5 pixels above the point

		r.Text(label, textX, textY)
	}
}

// LegendLeft is a legend that is designed for longer series lists.
func LegendLeft(currentChart *chart.Chart, userDefaults ...chart.Style) chart.Renderable {
	return func(render chart.Renderer, _ chart.Box, chartDefaults chart.Style) {
		legendDefaults := chart.Style{
			FillColor:   drawing.ColorWhite,
			FontColor:   chart.DefaultTextColor,
			FontSize:    8.0,
			StrokeColor: chart.DefaultAxisColor,
			StrokeWidth: chart.DefaultAxisLineWidth,
		}

		var legendStyle chart.Style
		if len(userDefaults) > 0 {
			legendStyle = userDefaults[0].InheritFrom(chartDefaults.InheritFrom(legendDefaults))
		} else {
			legendStyle = chartDefaults.InheritFrom(legendDefaults)
		}

		// DEFAULTS
		legendPadding := chart.Box{
			Top:    5,
			Left:   -10,
			Right:  5,
			Bottom: 5,
		}
		lineTextGap := 5
		lineLengthMinimum := 25

		var (
			labels = make([]string, 0, 1)
			lines  = make([]chart.Style, 0, 1)
		)

		for index, s := range currentChart.Series {
			if !s.GetStyle().Hidden {
				if _, isAnnotationSeries := s.(chart.AnnotationSeries); !isAnnotationSeries {
					labels = append(labels, s.GetName())
					lines = append(lines, s.GetStyle().InheritFrom(styleDefaultsSeries(currentChart, index)))
				}
			}
		}

		legend := chart.Box{
			Top:  5,
			Left: 5,
			// bottom and right will be sized by the legend content + relevant padding.
		}

		legendContent := chart.Box{
			Top:    legend.Top + legendPadding.Top,
			Left:   legend.Left + legendPadding.Left,
			Right:  legend.Left + legendPadding.Left,
			Bottom: legend.Top + legendPadding.Top,
		}

		legendStyle.GetTextOptions().WriteToRenderer(render)

		// measure
		var labelCount int

		for i := 0; i < len(labels); i++ {
			if len(labels[i]) > 0 {
				labelBox := render.MeasureText(labels[i])

				if labelCount > 0 {
					legendContent.Bottom += chart.DefaultMinimumTickVerticalSpacing
				}

				legendContent.Bottom += labelBox.Height()
				right := legendContent.Left + labelBox.Width() + lineTextGap + lineLengthMinimum
				legendContent.Right = chart.MaxInt(legendContent.Right, right)
				labelCount++
			}
		}

		legend = legend.Grow(legendContent)
		legend.Right = legendContent.Right + legendPadding.Right
		legend.Bottom = legendContent.Bottom + legendPadding.Bottom

		chart.Draw.Box(render, legend, legendStyle)

		legendStyle.GetTextOptions().WriteToRenderer(render)

		var (
			label       string
			yCursor     = legendContent.Top
			tx          = legendContent.Left
			legendCount = 0
		)

		for x := 0; x < len(labels); x++ {
			label = labels[x]
			if len(label) > 0 {
				if legendCount > 0 {
					yCursor += chart.DefaultMinimumTickVerticalSpacing
				}

				tb := render.MeasureText(label)

				ty := yCursor + tb.Height()
				render.Text(label, tx, ty)

				th2 := tb.Height() >> 1

				lx := tx + tb.Width() + lineTextGap
				ly := ty - th2
				lx2 := legendContent.Right - legendPadding.Right

				render.SetStrokeColor(lines[x].GetStrokeColor())
				render.SetStrokeWidth(lines[x].GetStrokeWidth())
				render.SetStrokeDashArray(lines[x].GetStrokeDashArray())

				render.MoveTo(lx, ly)
				render.LineTo(lx2, ly)
				render.Stroke()

				yCursor += tb.Height()
				legendCount++
			}
		}
	}
}

func styleDefaultsSeries(c *chart.Chart, seriesIndex int) chart.Style {
	return chart.Style{
		DotColor:    c.GetColorPalette().GetSeriesColor(seriesIndex),
		StrokeColor: c.GetColorPalette().GetSeriesColor(seriesIndex),
		StrokeWidth: chart.DefaultSeriesLineWidth,
		Font:        c.GetFont(),
		FontSize:    chart.DefaultFontSize,
	}
}
