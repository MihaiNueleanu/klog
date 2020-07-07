package cmd

import (
	"sort"
	"strings"
	"time"
)

// CutUrls - crop the length of urls
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

// SortByDate - orders the logs by request date
func SortByDate(logs []Log) []Log {
	sort.Slice(logs, func(i, j int) bool { 
		return logs[i].Time.Before(logs[j].Time)
	})
	
	return logs;
}

// FilterByStatus - filter logs by status code
func FilterByStatus(logs []Log, status int) []Log{
	var result = []Log{}
	for _, log := range logs {
		if(log.Status == status ) {
			result = append(result, log)
		}
	}

	return result
}

// FilterByLastNumberOfHours - filter by timestamp
func FilterByLastNumberOfHours(logs []Log, number int) []Log{
	now := time.Now()
	time := now.Add(-time.Hour * time.Duration(number))

	for i:=0; i< len(logs); i++ {
			if logs[i].Time.Before(time) {
				logs = append(logs[:i], logs[i+1:]...)
				i--
			}
	}

	return logs
}