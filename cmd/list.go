package cmd

import (
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
  Use:   "list",
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
		logs = CutUrls(logs)
	 
		RenderTable(logs)
  },
}

func init() {
  rootCmd.AddCommand(listCmd)
}
