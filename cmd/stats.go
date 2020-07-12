package cmd

import (
	"context"

	"github.com/mum4k/termdash"
	"github.com/mum4k/termdash/container"
	"github.com/mum4k/termdash/linestyle"
	"github.com/mum4k/termdash/terminal/termbox"
	"github.com/mum4k/termdash/terminal/terminalapi"
	"github.com/mum4k/termdash/widgets/barchart"
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
	values, labels := getData(logs)

	term, _ := termbox.New()

	defer term.Close()

	ctx, cancel := context.WithCancel(context.Background())
	bc, _ := barchart.New(
		barchart.ShowValues(),
		barchart.Labels(labels),
	)

	bc.Values(values, max(values) + 1)
	
	c, _ := container.New(
		term,
		container.Border(linestyle.Light),
		container.BorderTitle("PRESS Q TO QUIT"),
		container.PlaceWidget(bc),
	)
	

	quitter := func(k *terminalapi.Keyboard) {
		if k.Key == 'q' || k.Key == 'Q' {
			cancel()
		}
	}	

	err := termdash.Run(ctx, term, c, termdash.KeyboardSubscriber(quitter))

	if err != nil {
		panic(err)
	}
}

func getData(logs []Log) ([]int, []string) {
	startTime := logs[0].Time
	endTime := logs[len(logs)-1].Time
	firstDay := startTime.YearDay()
	lastDay := endTime.YearDay()
	numberOfPoints := lastDay-firstDay + 1

	values := make([]int, numberOfPoints)
	labels := make([]string, numberOfPoints)

	for _, log := range logs {
		day := log.Time.YearDay()
		index := day - firstDay

		values[index] = values[index] + 1	
		labels[index] = log.Time.Format("02/01");
	}

	return values, labels
}

func max(values []int) int {
	max := 0
	
	for _, val := range values {
		if (max < val) {
			max = val
		}
	}

	return max
}

func init() {
  rootCmd.AddCommand(statsCmd)
}
