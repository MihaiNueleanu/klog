package cmd

import (
	"sort"
 	"strings"
)

func CutUrls(lines []Log) []Log{
	for i, line := range lines {
		tmp := []rune(line.URL)
		line.URL = string(tmp[0:30])

		tmp2 := []rune(line.UserAgent)
		line.UserAgent = string(tmp2[0:20])

		lines[i] = line
	}


	return lines
}


// FilterLogsByFileTypes - filtering based on file type
func FilterLogsByFileTypes(lines []Log) []Log {
	types := []string {".css", ".png", ".svg", ".woff"}
	
	for i:=0; i< len(lines); i++ {
		for _, t := range types {
			if strings.Contains(lines[i].URL, t) {
				lines = append(lines[:i], lines[i+1:]...)
				i--
				break;
			}
		}
	}

	return lines
}

// SortByDate orders the logs by request date
func SortByDate(logs []Log) []Log {
	sort.Slice(logs, func(i, j int) bool { 
		return logs[i].Time.Before(logs[j].Time)
	})
	
	return logs;
}