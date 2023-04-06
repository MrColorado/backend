package scraper

// Scraper of each website should implement this interface
type Scraper interface {
	ScrapeNovel(novelName string)
	ScrapeNovelStart(novelName string, startChapter int)
	ScrapeNovelStartEnd(novelName string, startChapter int, endChapter int)
}
