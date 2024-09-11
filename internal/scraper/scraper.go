package scraper

import (
	"bufio"
	"context"
	"strings"
	"time"

	"github.com/gocolly/colly/v2"
	"github.com/yanilov/wc-scraper/internal/wordfilter"
)

type Scraper struct {
	config     ScraperConfig
	collector  *colly.Collector
	wordCount  map[string]int
	reduceOp   chan map[string]int
	reduceDone chan struct{}
}

func NewScraper(ctx context.Context, config ScraperConfig, filter wordfilter.WordFilter) *Scraper {

	s := &Scraper{
		config: config,
		collector: colly.NewCollector(
			colly.Async(),
			colly.IgnoreRobotsTxt(),
			colly.MaxDepth(0),
		),
		wordCount:  make(map[string]int),
		reduceOp:   make(chan map[string]int, config.Parallelism),
		reduceDone: make(chan struct{}),
	}

	// initialize the reduce step
	go func() {
		for wordCount := range s.reduceOp {
			for word, count := range wordCount {
				s.wordCount[word] += count
			}
		}
		// signal reduce is done
		s.reduceDone <- struct{}{}
		close(s.reduceDone)
	}()

	s.collector.SetRequestTimeout(2 * time.Second)

	s.collector.Limit(&colly.LimitRule{
		DomainGlob:  "*",
		Parallelism: config.Parallelism,
		RandomDelay: 23 * time.Millisecond,
		Delay:       1 * time.Second,
	})
	s.collector.OnHTML(config.Selector, func(e *colly.HTMLElement) {
		wordCount := map[string]int{}
		scanner := bufio.NewScanner(strings.NewReader(e.Text))
		scanner.Split(bufio.ScanWords)
		for scanner.Scan() {
			//non-blocking select, cancelling if the context is done
			select {
			case <-ctx.Done():
				return
			default:
				word := scanner.Text()
				if filter(word) {
					// safe to increment because indexing a map returns the zero falue if the key does not exist
					wordCount[word] += 1
				}
			}
		}
		s.reduceOp <- wordCount
	})

	return s
}

func (s *Scraper) Visit(url string) error {
	return s.collector.Visit(url)
}

func (s *Scraper) Wait() {
	s.collector.Wait()
	close(s.reduceOp)
	<-s.reduceDone
}

func (s *Scraper) WordCount() map[string]int {
	return s.wordCount
}
