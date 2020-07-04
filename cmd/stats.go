package cmd

import (
	"strings"
	"github.com/spf13/cobra"
	"fmt"
	"../logparser"
	"log"
	"github.com/olekukonko/tablewriter"
	"strconv"
	"os"
	"sort"
	"compress/gzip"
	"io/ioutil"
	"bufio"
	"time"
	"bytes"
	"regexp"
)


var statsCmd = &cobra.Command{
  Use:   "stats",
  Short: "Print the version number of KLog",
  Long:  `All software has versions. This is KLog's`,
  Run: func(cmd *cobra.Command, args []string) {
    parse()
  },
}


func readFiles() []string {
	files, err := ioutil.ReadDir("logs")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Print("Found files: ")
	fileNames := []string {}
	for _, f := range files {
		fmt.Print(f.Name() + ", ")
		fileNames = append(fileNames, f.Name())
	}
	fmt.Println("")
	fmt.Println("")

	return fileNames
}

func processArchives(fileNames []string) []string {
	var lines []string

	for _, f := range fileNames{
			file, err := os.Open(`logs/` + f)
			if err != nil {
				log.Fatal(err)
			}

			if strings.Contains(f, "gz") {
				reader, err := gzip.NewReader(file)

				if err != nil {
					log.Fatal(err)
				}

				scanner := bufio.NewScanner(reader)
				for scanner.Scan() {
					lines = append(lines, scanner.Text())
				}
			} else {
				scanner := bufio.NewScanner(file)
				for scanner.Scan() {
					lines = append(lines, scanner.Text())
				}
			}
	}

	return lines

}

func parse() {
	fileNames := readFiles()
	lines := processArchives(fileNames)

	logs := parseLines(lines)
	logs = filterLogs(logs)
	logs = cutUrls(logs)
 
	renderTable(logs)
}

func renderTable(lines []Line){
	
	table := tablewriter.NewWriter(os.Stdout)

	for _, line := range lines {
		fmt.Println(line)
		fmt.Printf("remote host: %s\n", line.RemoteHost)
		fmt.Printf("time: %s\n", line.Time)
		fmt.Printf("url: %s\n", line.URL)
		fmt.Printf("status: %d\n", line.Status)
		fmt.Println("------------------------------------------")


		table.Append([]string { 
			line.RemoteHost, 
			line.Time.Format("2006/01/02"),
			line.Time.Format("15:04:05"), 
			line.URL, 
			strconv.Itoa(line.Status),
		})
	}

	table.SetHeader([]string{"Source", "Day", "Time", "URL", "Code"})
	table.SetRowLine(true)
	table.SetAutoMergeCells(true)

	table.Render()
}

func filterLogs(lines []Line) []Line {
	for i:=0; i< len(lines); i++ {
		if strings.Contains(lines[i].URL, "css") ||
			strings.Contains(lines[i].URL, "png") ||
			strings.Contains(lines[i].URL, "svg") {
			 lines = append(lines[:i], lines[i+1:]...)
			 i--
		}
	}

	sort.Slice(lines, func(i, j int) bool { 
		return lines[i].Time.Before(lines[j].Time)
	})

	return lines
}


func cutUrls(lines []Line) []Line{
	for _, line := range lines {
		tmp := []rune(line.URL)
		line.URL = string(tmp[0:15])

		tmp2 := []rune(line.UserAgent)
		line.UserAgent = string(tmp2[0:15])

	}

	return lines
}


func parseLines(lines []string) []Line {
	var items []Line

	for _, line := range lines {
		var buffer bytes.Buffer
		buffer.WriteString(`^(\S+)\s`)                  // 1) IP
		buffer.WriteString(`\S+\s+`)                    // remote logname
		buffer.WriteString(`(?:\S+\s+)+`)               // remote user
		buffer.WriteString(`\[([^]]+)\]\s`)             // 2) date
		buffer.WriteString(`"(\S*)\s?`)                 // 3) method
		buffer.WriteString(`(?:((?:[^"]*(?:\\")?)*)\s`) // 4) URL
		buffer.WriteString(`([^"]*)"\s|`)               // 5) protocol
		buffer.WriteString(`((?:[^"]*(?:\\")?)*)"\s)`)  // 6) or, possibly URL with no protocol
		buffer.WriteString(`(\S+)\s`)                   // 7) status code
		buffer.WriteString(`(\S+)\s`)                   // 8) bytes
		buffer.WriteString(`"((?:[^"]*(?:\\")?)*)"\s`)  // 9) referrer
		buffer.WriteString(`"(.*)"$`)                   // 10) user agent

		re1, err := regexp.Compile(buffer.String())
		if err != nil {
			log.Fatalf("regexp: %s", err)
		}
		result := re1.FindStringSubmatch(line)

		lineItem := new(Line)
		lineItem.RemoteHost = result[1]
		// [05/Oct/2014:04:06:21 -0500]
		value := result[2]
		layout := "02/Jan/2006:15:04:05 -0700"
		t, _ := time.Parse(layout, value)
		lineItem.Time = t
		lineItem.Request = result[3] + " " + result[4] + " " + result[5]
		status, err := strconv.Atoi(result[7])
		if err != nil {
			status = 0
		}
		bytes, err := strconv.Atoi(result[8])
		if err != nil {
			bytes = 0
		}
		lineItem.Status = status
		lineItem.Bytes = bytes
		lineItem.Referer = result[9]
		lineItem.UserAgent = result[10]
		url := result[4]
		altURL := result[6]
		if url == "" && altURL != "" {
			url = altURL
		}
		lineItem.URL = url
		items = append(items, *lineItem)
		//for k, v := range result {
		//	fmt.Printf("%d. %s\n", k, v)
		//}
	}
	return items
}

func parseLogs() {
	lines, err := logparser.Parse("access.log")

	if err != nil {
		log.Fatal(err)
	}

	for i:=0; i< len(lines); i++ {
		if !strings.Contains(lines[i].URL, "css") ||
			!strings.Contains(lines[i].URL, "png") ||
			!strings.Contains(lines[i].URL, "svg") {
			 lines = append(lines[:i], lines[i+1:]...)
			 i--
		}
	}

	sort.Slice(lines, func(i, j int) bool { 
		return lines[i].Time.Before(lines[j].Time)
	})

	table := tablewriter.NewWriter(os.Stdout)

	for _, line := range lines {
		//fmt.Println(line)
		fmt.Printf("remote host: %s\n", line.RemoteHost)
		fmt.Printf("time: %s\n", line.Time)
		fmt.Printf("url: %s\n", line.URL)
		fmt.Printf("status: %d\n", line.Status)
		fmt.Println("------------------------------------------")


		table.Append([]string { 
			line.RemoteHost, 
			line.Time.Format("2006/01/02"),
			line.Time.Format("15:04:05"), 
			line.URL, 
			strconv.Itoa(line.Status),
		})
	}

	table.SetHeader([]string{"Source", "Day", "Time", "URL", "Code"})
	table.SetRowLine(true)
	table.SetAutoMergeCells(true)

	table.Render()

}

type Line struct {
	RemoteHost string
	Time       time.Time
	Request    string
	Status     int
	Bytes      int
	Referer    string
	UserAgent  string
	URL        string
}

func init() {
  rootCmd.AddCommand(statsCmd)
}
