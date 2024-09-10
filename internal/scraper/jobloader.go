package scraper

import (
	"bufio"
	"io"
	"net/http"
	"os"
)

func LoadJobsFromUrls(mainUrl string, lineBuffer int) (<-chan ScrapeJobSpec, error) {
	resp, err := http.Get(mainUrl)
	if err != nil {
		return nil, err
	}

	result := loadJobsFromReaderCloser(resp.Body, lineBuffer)
	return result, nil
}

func LoadFromFile(filePath string, lineBuffer int) (<-chan ScrapeJobSpec, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}

	result := loadJobsFromReaderCloser(file, lineBuffer)
	return result, nil
}

func loadJobsFromReaderCloser(reader io.ReadCloser, lineBuffer int) <-chan ScrapeJobSpec {
	result := make(chan ScrapeJobSpec, lineBuffer)
	scanner := bufio.NewScanner(reader)
	scanner.Split(bufio.ScanLines)
	// start a goroutine to read the scanner and send the urls to the result channel, line by line
	go func() {
		defer close(result)
		defer reader.Close()
		for scanner.Scan() {
			url := scanner.Text()
			result <- ScrapeJobSpec{url: url}
		}
		// if there was an error reading the scanner, send the error to the result channel
		if err := scanner.Err(); err != nil {
			result <- ScrapeJobSpec{err: err}
		}
	}()

	return result
}
