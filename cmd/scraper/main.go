package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/yanilov/wc-scraper/internal/analytics"
	"github.com/yanilov/wc-scraper/internal/bank"
	"github.com/yanilov/wc-scraper/internal/scraper"
	"github.com/yanilov/wc-scraper/internal/wordfilter"
)

const (
	bankUrl = "https://raw.githubusercontent.com/dwyl/english-words/master/words.txt"
	jobsUrl = "https://drive.usercontent.google.com/u/0/uc?id=1TF4RPuj8iFwpa-lyhxG67V8NDlktmTGi&export=download"

	// safety parameter to avoid loading too many pages and getting blocked.
	// set it to 0 to load all pages
	loadJobPageCutoff   = 3
	loadJobBackpressure = 10

	scrapeParallelism = 6
	scrapeSelector    = "article p, article h1, article h2, article h3, article h4, article h5, article h6"

	topK = 10
)

func main() {

	ctx := buildSignalContext()

	alphaAndMinLenFilter := wordfilter.NewAggregateFilter(
		wordfilter.NewMinLengthFilter(3),
		wordfilter.NewAlphaOnlyFilter(),
	)

	bankFuture, err := buildWordBank(ctx, bankUrl, alphaAndMinLenFilter)
	if err != nil {
		panic(err)
	}
	jobStream, err := scraper.LoadJobsFromUrls(ctx, jobsUrl, scraper.ScrapeJobLoaderConfig{
		Backpressure: loadJobBackpressure,
		PageCutoff:   loadJobPageCutoff,
	})
	if err != nil {
		panic(err)
	}

	bank, ok := <-bankFuture
	if !ok {
		panic("future did not resolve to a word bank")
	}

	// create an aggregate filter over alphas, min length and the word bank
	scraperWordfilter := wordfilter.NewAggregateFilter(
		alphaAndMinLenFilter,
		wordfilter.NewWordBankFilter(bank),
	)

	scrapeConfig := scraper.ScraperConfig{
		Parallelism: scrapeParallelism,
		Selector:    scrapeSelector,
	}
	scraper := scraper.NewScraper(ctx, scrapeConfig, scraperWordfilter)

	report := NewErrorReport()
mainloop:
	for {
		select {
		case job, ok := <-jobStream:
			// if all jobs have been processed
			if !ok {
				break mainloop
			}
			url, err := job.Unpack()
			// append to report if there was an error loading the job
			if err != nil {
				report.LoadUrlsErrors = append(report.LoadUrlsErrors, err)
				continue
			}
			scraper.Visit(url)

		case <-ctx.Done():
			return
		}
	}

	// wait for all scraping to finish
	scraper.Wait()

	// print the top K words
	wc := scraper.WordCount()
	wcTopK := analytics.TopK(wc, topK)
	bytes, err := json.MarshalIndent(wcTopK, "", "  ")
	if err != nil {
		panic(err)
	}
	fmt.Fprintln(os.Stdout, string(bytes))

	// print to stderr on any errors
	if !report.IsEmpty() {
		bytes, err := json.MarshalIndent(report.Error(), "", "  ")
		if err != nil {
			panic(err)
		}
		fmt.Fprintln(os.Stderr, string(bytes))
	}
}

func buildWordBank(ctx context.Context, bankUrl string, filter wordfilter.WordFilter) (<-chan bank.WordBank, error) {
	bankFuture, err := scraper.LoadBankFromUrl(ctx, bankUrl, filter)
	return bankFuture, err
}

// buildSignalContext creates a context that is cancelled when the process receives a SIGINT or SIGTERM signal
func buildSignalContext() context.Context {
	ctx, cancel := context.WithCancel(context.Background())

	c := make(chan os.Signal, 1)

	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	// on ctrl+c or ctrl+d, cancel the context
	go func() {
		select {
		case <-c:
			cancel()
			fmt.Fprintln(os.Stderr, "cancelled by user")
		case <-ctx.Done():
		}
	}()
	return ctx
}

type errorReport struct {
	LoadUrlsErrors  []error
	UrlScrapeErrors map[string]error
}

// assert that ErrorReport implements error
var _ error = (*errorReport)(nil)

// Error implements error.
func (e *errorReport) Error() string {
	return fmt.Sprintf("errored load jobs: %v, errored url scrapes: %v", e.LoadUrlsErrors, e.UrlScrapeErrors)
}

func (e *errorReport) IsEmpty() bool {
	return len(e.LoadUrlsErrors) == 0 && len(e.UrlScrapeErrors) == 0
}

func NewErrorReport() *errorReport {
	return &errorReport{
		LoadUrlsErrors:  make([]error, 0),
		UrlScrapeErrors: make(map[string]error),
	}
}
