package cmd

import (
	"log"

	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"

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

	if err := ui.Init(); err != nil {
		log.Fatalf("failed to initialize termui: %v", err)
	}

	defer ui.Close()

	chart :=widgets.NewBarChart()

	chart.Data = values
	chart.Labels = labels
	chart.Title = "Plot chart"
	chart.SetRect(0, 0, 100, 20)
	chart.BarWidth = 7
	chart.BarColors = []ui.Color{ui.ColorRed, ui.ColorGreen}
	chart.LabelStyles = []ui.Style{ui.NewStyle(ui.ColorBlue)}
	chart.NumStyles = []ui.Style{ui.NewStyle(ui.ColorYellow)}

	ui.Render(chart)

	uiEvents := ui.PollEvents()
	for {
		e := <-uiEvents
		switch e.ID {
		case "q", "<C-c>":
						return
		}
	}

}

func getData(logs []Log) ([]float64, []string) {
	startTime := logs[0].Time
	endTime := logs[len(logs)-1].Time
	firstDay := startTime.YearDay()
	lastDay := endTime.YearDay()
	numberOfPoints := lastDay-firstDay + 1

	values := make([]float64, numberOfPoints)
	labels := make([]string, numberOfPoints)

	for _, log := range logs {
		
		day := log.Time.YearDay()
		index := day - firstDay

		values[index] = values[index] + 1	
		labels[index] = log.Time.Format("01/02");
	}

	return values, labels
}

func init() {
  rootCmd.AddCommand(statsCmd)
}
