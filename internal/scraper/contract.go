package scraper

type ScrapeJobSpec struct {
	url string
	err error
}

func (s *ScrapeJobSpec) Unpack() (string, error) {
	return s.url, s.err
}

type ScrapeJobLoaderConfig struct {
	// cut off the number of pages to scrape, similar to head in unix. useful for testing and avoiding getting blocked
	PageCutoff int
	// the number of jobs to buffer in the output channel, to avoid pressuring the scraper
	Backpressure int
}
