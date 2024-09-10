package scraper

import (
	"bufio"
	"context"
	"io"
	"net/http"

	"github.com/yanilov/wc-scraper/internal/bank"
	"github.com/yanilov/wc-scraper/internal/wordfilter"
)

func LoadBankFromUrl(ctx context.Context, url string, filter wordfilter.WordFilter) (<-chan bank.WordBank, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	result := loadBankFromReaderCloser(ctx, resp.Body, filter)
	return result, nil
}

func loadBankFromReaderCloser(ctx context.Context, reader io.ReadCloser, filter wordfilter.WordFilter) <-chan bank.WordBank {

	scanner := bufio.NewScanner(reader)
	scanner.Split(bufio.ScanWords)

	result := make(chan bank.WordBank)

	go func() {
		defer close(result)
		defer reader.Close()

		bank := bank.New()
		for scanner.Scan() {
			//non-blocking select, cancelling if the context is done
			select {
			case <-ctx.Done():
				return
			default:
				word := scanner.Text()
				if filter(word) {
					bank.Add(word)
				}
			}
		}
		result <- bank
	}()

	return result
}
