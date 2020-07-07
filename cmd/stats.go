package cmd

import (
	"context"
	"fmt"

	"github.com/mum4k/termdash"
	"github.com/mum4k/termdash/cell"
	"github.com/mum4k/termdash/container"
	"github.com/mum4k/termdash/linestyle"
	"github.com/mum4k/termdash/terminal/termbox"
	"github.com/mum4k/termdash/terminal/terminalapi"
	"github.com/mum4k/termdash/widgets/linechart"
	"github.com/spf13/cobra"
)

var statsCmd = &cobra.Command{
  Use:   "stats",
  Short: "Print the version number of KLog",
  Long:  `All software has versions. This is KLog's`,
  Run: func(cmd *cobra.Command, args []string) {
		fileNames := ReadFiles()
		lines := LoadFiles(fileNames)

		logs := ParseLines(lines)
		logs = FilterLogsByFileTypes(logs)
		
		if(StatusCode != 0) {
			logs = FilterByStatus(logs, StatusCode)			
		}

		if(HoursBack != 0){
			logs = FilterByLastNumberOfHours(logs, HoursBack)
		}

		logs = SortByDate(logs)

		plot(logs)
	 
  },
}

func plot(logs []Log) {
	term, err := termbox.New()
	if err != nil {
		panic(err)
	}

	fmt.Println("reached term close")
	
	defer term.Close()

	ctx, cancel := context.WithCancel(context.Background())
	lc, err := linechart.New(
		linechart.AxesCellOpts(cell.FgColor(cell.ColorRed)),
		linechart.YLabelCellOpts(cell.FgColor(cell.ColorGreen)),
		linechart.XLabelCellOpts(cell.FgColor(cell.ColorCyan)),
	)

	if err != nil {
		panic(err)
	}

	go playLineChart(ctx, lc, logs)

	box, err := container.New(
		term,
		container.Border(linestyle.Light),
		container.BorderTitle("PRESS Q TO QUIT"),
		container.PlaceWidget(lc),
	)

	if err != nil {
		panic(err)
	}

	quitter := func(k *terminalapi.Keyboard) {
		if k.Key == 'q' || k.Key == 'Q' {
			cancel()
		}
	}

	if err := termdash.Run(ctx, term, box, termdash.KeyboardSubscriber(quitter)); err != nil {
		panic(err)
	}
}

// playLineChart continuously adds values to the LineChart, once every delay.
// Exits when the context expires.
func playLineChart(ctx context.Context, lc *linechart.LineChart, logs []Log) {
	startTime := logs[0].Time
	endTime := logs[len(logs)-1].Time
	numberOfPoints := int(endTime.Sub(startTime).Hours()) + 1

	values := make([]float64, numberOfPoints)
	labels := make(map[int]string, numberOfPoints)

	for k, log := range logs {
		hour := log.Time.Sub(startTime).Hours()
		index := int(hour)

		values[index] = values[index] + 1	
		labels[k] = log.Time.Format("2006/01/02 15:04:05");
	}

	err := lc.Series(
		"first", 
		values,
		linechart.SeriesCellOpts(cell.FgColor(cell.ColorBlue)),
		linechart.SeriesXLabels(labels),
	)

	if err != nil {
		panic(err)
	}
}

func init() {
  rootCmd.AddCommand(statsCmd)
}
