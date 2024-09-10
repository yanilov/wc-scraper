package scraper

type ScrapeJobSpec struct {
	url string
	err error
}

func (s *ScrapeJobSpec) Unpack() (string, error) {
	return s.url, s.err
}
