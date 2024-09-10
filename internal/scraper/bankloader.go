package scraper

import (
	"bufio"
	"io"
	"net/http"

	"github.com/yanilov/wc-scraper/internal/bank"
	"github.com/yanilov/wc-scraper/internal/wordfilter"
)

func LoadBankFromUrl(url string, filter wordfilter.WordFilter) (<-chan bank.WordBank, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	result := loadBankFromReaderCloser(resp.Body, filter)
	return result, nil
}

func loadBankFromReaderCloser(reader io.ReadCloser, filter wordfilter.WordFilter) <-chan bank.WordBank {

	scanner := bufio.NewScanner(reader)
	scanner.Split(bufio.ScanWords)

	result := make(chan bank.WordBank)

	go func() {
		defer close(result)
		defer reader.Close()

		bank := bank.New()
		for scanner.Scan() {
			word := scanner.Text()
			if filter(word) {
				bank.Add(word)
			}
		}
	}()

	return result
}
