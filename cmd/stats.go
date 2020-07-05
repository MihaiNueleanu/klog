package cmd

import (
	"github.com/spf13/cobra"
)

var statusCode int
var hoursBack int

var listCmd = &cobra.Command{
  Use:   "list",
  Short: "Print the version number of KLog",
  Long:  `All software has versions. This is KLog's`,
  Run: func(cmd *cobra.Command, args []string) {
		fileNames := ReadFiles()
		lines := LoadFiles(fileNames)

		logs := ParseLines(lines)
		logs = FilterLogsByFileTypes(logs)
		
		if(statusCode != 0) {
			logs = FilterByStatus(logs, statusCode)			
		}

		if(hoursBack != 0){
			logs = FilterByLastNumberOfHours(logs, hoursBack)
		}

		logs = SortByDate(logs)
		logs = CutUrls(logs)
	 
		RenderTable(logs)
  },
}

func init() {
	listCmd.Flags().IntVarP(&statusCode, "code", "c", 200, "Status code to filter by")
	listCmd.Flags().IntVarP(&hoursBack, "back", "b", 168, "Number of hours back")

  rootCmd.AddCommand(listCmd)
}
