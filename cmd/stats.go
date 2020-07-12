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

	chart :=widgets.NewPlot()

	chart.DataLabels = labels
	chart.Data = make([][]float64, 1)
	chart.Data[0] = values
	chart.Title = "Plot chart"
	chart.SetRect(3, 3, 100, 20)
	chart.AxesColor = ui.ColorWhite
	chart.LineColors[0] = ui.ColorGreen
	chart.ShowAxes = true

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
	numberOfPoints := int(endTime.Sub(startTime).Hours()) + 1

	values := make([]float64, numberOfPoints)
	labels := make([]string, numberOfPoints)

	for k, log := range logs {
		hour := log.Time.Sub(startTime).Hours()
		index := int(hour)

		values[index] = values[index] + 1	
		labels[k] = log.Time.Format("01/02 15:04");
	}

	return values, labels
}

func init() {
  rootCmd.AddCommand(statsCmd)
}
