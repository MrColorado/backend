package scraper

// Scraper of each website should implement this interface
type Scraper interface {
	ScrapeNovel(novelName string)
}
