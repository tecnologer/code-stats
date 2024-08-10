package chartui

import (
	"fmt"
	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
	"github.com/tecnologer/code-stats/pkg/models"
	"syscall"
	"unsafe"
)

func Draw(stats *models.StatsCollection, statType models.StatType, isDiff bool, languages ...string) error {
	if err := ui.Init(); err != nil {
		return fmt.Errorf("failed to initialize termui: %w", err)
	}

	defer ui.Close()

	// Sample data points
	series := make(map[string][]float64)

	for _, key := range stats.KeysSorted() {
		languageStats := stats.Get(key)

		for _, stat := range languageStats {
			if !stat.IsInLanguageList(languages) {
				continue
			}

			if serie, ok := series[stat.Name]; ok {
				serie = append(serie, stat.ValueOf(statType))
				continue
			}

			series[stat.Name] = []float64{stat.ValueOf(statType)}
		}
	}

	// Set up the line chart
	lc := widgets.NewPlot()
	lc.Title = fmt.Sprintf("%s Count", statType.Title())

	width := 0.5 * getWidth()
	if width < 150 {
		width = 150
	}

	lc.SetRect(0, 0, int(width), 250)
	lc.AxesColor = ui.ColorWhite
	lc.Marker = widgets.MarkerBraille

	for label, serie := range series {
		lc.Data = append(lc.Data, serie)
		lc.DataLabels = append(lc.DataLabels, label)
	}

	// Rendering the line chart
	ui.Render(lc)

	// Event loop to keep the UI running
	uiEvents := ui.PollEvents()
	for {
		e := <-uiEvents
		if e.ID == "q" || e.Type == ui.KeyboardEvent {
			break
		}
	}

	return nil
}

type winsize struct {
	Row    uint16
	Col    uint16
	Xpixel uint16
	Ypixel uint16
}

func getWidth() float32 {
	ws := &winsize{}
	retCode, _, errno := syscall.Syscall(syscall.SYS_IOCTL,
		uintptr(syscall.Stdin),
		uintptr(syscall.TIOCGWINSZ),
		uintptr(unsafe.Pointer(ws)))

	if int(retCode) == -1 {
		panic(errno)
	}
	return float32(ws.Col)
}
