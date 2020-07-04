package cmd

import (
  "fmt"
  "github.com/spf13/cobra"
)


var versionCmd = &cobra.Command{
  Use:   "version",
  Short: "Print the version number of KLog",
  Long:  `All software has versions. This is KLog's`,
  Run: func(cmd *cobra.Command, args []string) {
    fmt.Println("KLog server log analyzer v0.1")
  },
}


func init() {
  rootCmd.AddCommand(versionCmd)
}
