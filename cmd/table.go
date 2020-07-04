package cmd

import (
	"github.com/olekukonko/tablewriter"
	"strconv"
	"os"
)

// RenderTable for pretty logs 
func RenderTable(lines []Log){
	
	table := tablewriter.NewWriter(os.Stdout)

	for _, line := range lines {
		table.Append([]string { 
			line.RemoteHost, 
			line.Time.Format("2006/01/02"),
			line.Time.Format("15:04:05"), 
			line.URL, 
			strconv.Itoa(line.Status),
		})
	}

	table.SetHeader([]string{"Source", "Day", "Time", "URL", "Code"})
	// table.SetRowLine(true)
	// table.SetAutoMergeCells(true)

	table.Render()
}

