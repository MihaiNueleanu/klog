package cmd

import (
	"github.com/spf13/cobra"
)

var statsCmd = &cobra.Command{
  Use:   "stats",
  Short: "Print the version number of KLog",
  Long:  `All software has versions. This is KLog's`,
  Run: func(cmd *cobra.Command, args []string) {
    raw()
  },
}

func raw() {
	fileNames := ReadFiles()
	lines := LoadFiles(fileNames)

	logs := ParseLines(lines)
	logs = FilterLogsByFileTypes(logs)
	logs = SortByDate(logs)
	logs = CutUrls(logs)
 
	RenderTable(logs)
}

func init() {
  rootCmd.AddCommand(statsCmd)
}
