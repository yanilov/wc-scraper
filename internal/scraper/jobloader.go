package scraper

import (
	"bufio"
	"context"
	"io"
	"net/http"
	"os"
)

func LoadJobsFromUrls(ctx context.Context, mainUrl string, config ScrapeJobLoaderConfig) (<-chan ScrapeJobSpec, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", mainUrl, nil)
	if err != nil {
		return nil, err
	}
	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		return nil, err
	}

	result := loadJobsFromReaderCloser(ctx, resp.Body, config)
	return result, nil
}

func LoadFromFile(ctx context.Context, filePath string, config ScrapeJobLoaderConfig) (<-chan ScrapeJobSpec, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}

	result := loadJobsFromReaderCloser(ctx, file, config)
	return result, nil
}

func loadJobsFromReaderCloser(ctx context.Context, reader io.ReadCloser, config ScrapeJobLoaderConfig) <-chan ScrapeJobSpec {
	result := make(chan ScrapeJobSpec, config.Backpressure)
	scanner := bufio.NewScanner(reader)
	scanner.Split(bufio.ScanLines)
	// start a goroutine to read the scanner and send the urls to the result channel, line by line
	go func() {
		defer close(result)
		defer reader.Close()

		for i := 0; scanner.Scan(); i++ {
			//break early if page cutoff is set
			if config.PageCutoff > 0 && i >= config.PageCutoff {
				break
			}
			//non-blocking select, cancelling if the context is done
			select {
			case <-ctx.Done():
				return
			default:
				url := scanner.Text()
				result <- ScrapeJobSpec{url: url}
			}
		}
		// if there was an error reading the scanner, send the error to the result channel
		if err := scanner.Err(); err != nil {
			result <- ScrapeJobSpec{err: err}
		}
	}()

	return result
}
