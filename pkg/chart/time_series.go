package chart

import (
	"fmt"
	"github.com/tecnologer/code-stats/ui"
	"github.com/wcharczuk/go-chart/v2"
)

type TimeSeries struct {
	chart.TimeSeries
}

// Render renders the series to the chart.
func (ats *TimeSeries) Render(renderer chart.Renderer, canvasBox chart.Box, xRange, yRange chart.Range, defaultStyle chart.Style) {
	style := ats.Style.InheritFrom(defaultStyle)

	// Render the line and dots as usual
	ats.TimeSeries.Render(renderer, canvasBox, xRange, yRange, style)

	// Set font and color for text rendering
	renderer.SetFont(style.GetFont())
	renderer.SetFontColor(chart.ColorAlternateGreen) // Set a bright color for visibility

	// X is filled with time.Time values
	for index, date := range ats.XValues {
		// Convert time.Time to float64 for the x-axis
		xValue := xRange.Translate(chart.TimeToFloat64(date))
		yValue := yRange.Translate(ats.YValues[index])

		// Format the label text
		label := fmt.Sprintf("%.2f", ats.YValues[index])

		// Measure text size
		textDimensions := renderer.MeasureText(label)

		// Calculate text position to ensure visibility
		textX := canvasBox.Left + xValue - textDimensions.Width()/2
		textY := canvasBox.Top + yValue - textDimensions.Height() - 10 // Adjusted for visibility

		// Debug output
		ui.Debugf("Text '%s' at (%d, %d), within canvas (%d, %d, %d, %d)",
			label, textX, textY, canvasBox.Left, canvasBox.Top, canvasBox.Right, canvasBox.Bottom)

		// Draw the text
		renderer.Text(label, textX, textY)
	}
}
