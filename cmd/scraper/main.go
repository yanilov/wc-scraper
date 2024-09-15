package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/alecthomas/kong"
	"github.com/yanilov/wc-scraper/internal/analytics"
	"github.com/yanilov/wc-scraper/internal/bank"
	"github.com/yanilov/wc-scraper/internal/config"
	"github.com/yanilov/wc-scraper/internal/scraper"
	"github.com/yanilov/wc-scraper/internal/wordfilter"
)

var stderrLogger *slog.Logger

func initLogger(level string) {
	var slogLevel slog.Level
	switch level {
	case "", "info", "INFO":
		slogLevel = slog.LevelInfo
	case "debug", "DEBUG":
		slogLevel = slog.LevelDebug
	case "warn", "WARN":
		slogLevel = slog.LevelWarn
	case "error", "ERROR":
		slogLevel = slog.LevelError
	default:
		panic("invalid log level")
	}

	stderrLogger = slog.New(slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{
		Level: slogLevel,
	}))
}

func init() {
	initLogger(os.Getenv("GO_LOG"))
}

var defaultConfig = config.AppConfig{
	ScraperConfig: scraper.ScraperConfig{
		Parallelism: 6,
		Selector:    "article p, article h1, article h2, article h3, article h4, article h5, article h6",
	},
	JobLoaderConfig: scraper.ScrapeJobLoaderConfig{
		Backpressure: 10,
		// safety parameter to avoid loading too many pages and getting blocked.
		// set it to 0 to load all pages
		PageCutoff: 3,
	},
	WordBankUrl: "https://raw.githubusercontent.com/dwyl/english-words/master/words.txt",
	JobsUrl:     "https://drive.usercontent.google.com/u/0/uc?id=1TF4RPuj8iFwpa-lyhxG67V8NDlktmTGi&export=download",
}

type CLI struct {
	Level  string `default:"info" help:"Log level."`
	Config string `default:"" optional:"" help:"Path to the config file."`
	Head   int    `arg:"" default:"10" help:"Number of top words to display."`
}

func (tk *CLI) Run(cfg *config.AppConfig, head int) error {
	// create a context that is cancelled when the process receives a SIGINT or SIGTERM signal
	ctx := buildSignalContext()

	// run the actual scraping
	res, err := scrapeTopK(ctx, cfg, head)
	if err != nil {
		return err
	}

	// print the top K words to stdout
	bytes, err := json.MarshalIndent(res, "", "  ")
	if err != nil {
		return err
	}

	fmt.Fprintln(os.Stdout, string(bytes))
	return nil
}

func main() {
	//parse cli args
	var cli CLI
	kctx := kong.Parse(&cli)

	var appconfig *config.AppConfig
	var err error

	switch cli.Config {
	case "":
		appconfig = &defaultConfig
		stderrLogger.Info("using config",
			"type", "default",
			"content", defaultConfig)
	default:
		appconfig, err = config.FromFile(cli.Config)
		if err != nil {
			panic(err)
		}
		stderrLogger.Info("using config",
			"type", "file",
			"content", appconfig)
	}

	err = kctx.Run(appconfig, cli.Head)
	if err != nil {
		panic(err)
	}
}

type wordCount map[string]int

func scrapeTopK(ctx context.Context, config *config.AppConfig, head int) (wordCount, error) {
	if head <= 0 {
		return nil, errors.New("head must be greater than 0")
	}
	if config.JobLoaderConfig.PageCutoff > 0 {
		stderrLogger.Info("loading partial pages",
			"page_cutoff", config.JobLoaderConfig.PageCutoff)
	}

	alphaAndMinLenFilter := wordfilter.NewAggregateFilter(
		wordfilter.NewMinLengthFilter(3),
		wordfilter.NewAlphaOnlyFilter(),
	)

	bankFuture, err := buildWordBank(ctx, config.WordBankUrl, alphaAndMinLenFilter)
	if err != nil {
		return nil, err
	}
	jobStream, err := scraper.LoadJobsFromUrls(ctx, config.JobsUrl, config.JobLoaderConfig)
	if err != nil {
		return nil, err
	}

	bank, ok := <-bankFuture
	if !ok {
		return nil, errors.New("future did not resolve to a word bank")
	}

	// create an aggregate filter over alphas, min length and the word bank
	scraperWordfilter := wordfilter.NewAggregateFilter(
		alphaAndMinLenFilter,
		wordfilter.NewWordBankFilter(bank),
	)

	scraper := scraper.NewScraper(ctx, config.ScraperConfig, scraperWordfilter)

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
			stderrLogger.Debug("scraping url",
				"url", url)
			scraper.Visit(url)

		case <-ctx.Done():
			return nil, ctx.Err()
		}
	}

	// wait for all scraping to finish
	scraper.Wait()

	// print the top K words
	wc := scraper.WordCount()
	wcTopK := analytics.TopK(wc, head)

	if !report.IsEmpty() {
		return wcTopK, report
	} else {
		return wcTopK, nil
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
			stderrLogger.Info("cancelled by user")
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
	bytes, err := json.MarshalIndent(*e, "", "  ")
	if err != nil {
		panic("must not fail to marshal error report")
	}
	return string(bytes)
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
