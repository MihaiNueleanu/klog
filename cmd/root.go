package cmd

import (
	"github.com/spf13/cobra"
)

var Source string

var rootCmd = &cobra.Command{
  Use:   "klog",
  Short: "Klog is a log analysing interface",
  Long: `A fast and simple way to parse through your server logs`,
}

// Execute executes the root command.
func Execute() error {
  return rootCmd.Execute()
}

func init() {
	rootCmd.PersistentFlags().StringP("author", "a", "Mihai Nueleanu", "author name for copyright attribution")
}

