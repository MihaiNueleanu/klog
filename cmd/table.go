package cmd

import (
	"os"
	"strconv"

	"github.com/mssola/user_agent"
	"github.com/olekukonko/tablewriter"
)

// RenderTable for pretty logs
func RenderTable(lines []Log){
	
	table := tablewriter.NewWriter(os.Stdout)
	
	table.SetHeader([]string{"Source", "Day", "Time", "URL", "Code", "Browser"})
	ua := user_agent.New("Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.11 (KHTML, like Gecko) Chrome/23.0.1271.97 Safari/537.11")


	for _, line := range lines {
		browser, _ := ua.Browser()

		table.Append([]string { 
			line.RemoteHost, 
			line.Time.Format("2006/01/02"),
			line.Time.Format("15:04:05"), 
			// line.Method,  
			line.URL, 
			strconv.Itoa(line.Status),
			browser,
		})
	}

	table.Render()
}

