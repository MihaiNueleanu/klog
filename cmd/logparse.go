package cmd

import (
	"bufio"
	"bytes"
	"compress/gzip"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// Log type
type Log struct {
	RemoteHost string
	Time       time.Time
	Request    string
	Method     string
	Status     int
	Bytes      int
	Referer    string
	UserAgent  string
	URL        string
}

// ParseLines - takes lines and return Log objects
func ParseLines(lines []string) []Log {
	var items []Log

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

		lineItem := new(Log)
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
		lineItem.Method = result[3]
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

// ReadFiles - discovers the files in the logs folder
func ReadFiles() []string {
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

// LoadFiles - loads the contents of the files in memory
func LoadFiles(fileNames []string) []string {
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