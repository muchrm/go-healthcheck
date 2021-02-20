package utils

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/muchrm/go-healthcheck/config"
)

// ReadCSVAndGetServerUrls receive csv file name and get server url from csv
func ReadCSVAndGetServerUrls(fileName string) ([]string, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return []string{}, errors.New(config.CSVNotFoundError)
	}
	urls, err := GetServerUrls(file)
	if err != nil {
		return []string{}, errors.New(config.CSVIncorrectFormatError)
	}
	return urls, nil
}

// GetServerUrls receive buffer file and check server url, if column not match url scheme it will be ignore
func GetServerUrls(reader io.Reader) ([]string, error) {
	r := csv.NewReader(reader)
	lines, err := r.ReadAll()
	if err != nil {
		return nil, err
	}
	URLs := []string{}
	for _, column := range lines {
		URLString := column[0]
		URLString = strings.TrimSpace(URLString)
		_, err := url.ParseRequestURI(URLString)
		if err != nil {
			fmt.Println("column url is incorrect format")
		} else {
			URLs = append(URLs, URLString)
		}

	}
	return URLs, err
}

// GetDurationString receive duration and return time to proceed (ms/sec/minutes) in string
func GetDurationString(t time.Duration) string {
	minutes := int(time.Duration(t % time.Hour).Minutes())
	seconds := int(time.Duration(t % time.Minute).Seconds())
	milliseconds := time.Duration(t % time.Second).Milliseconds()
	timeString := ""
	if minutes > 0 {
		timeString = fmt.Sprintf("%d minutes ", minutes)
	}
	if seconds > 0 {
		timeString = fmt.Sprintf("%s%d sec ", timeString, seconds)
	}
	if milliseconds > 0 {
		timeString = fmt.Sprintf("%s%d ms", timeString, milliseconds)
	}
	return timeString
}
