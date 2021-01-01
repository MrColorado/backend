package scraper

import (
	"github.com/gocolly/colly"
)

// Scraper of each website should implement this interface
type Scraper interface {
	ScrapeNovel(c *colly.Collector, novelName string, outputPath string)
	ScrapPartialNovel(c *colly.Collector, novelName string, startChapter int, endChapter int, outputPath string)
}
