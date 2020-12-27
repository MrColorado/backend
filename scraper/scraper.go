package scraper

import (
	"github.com/gocolly/colly"
)

// NovelData contain data of a chapter
type NovelData struct {
	Title     string
	Chapter   int
	Paragraph []string
}

// Scraper of each website should implement this interface
type Scraper interface {
	ScrapeNovel(c *colly.Collector, novelName string)
	ScrapPartialNovel(c *colly.Collector, novelName string, startChapter int, endChapter int)
}
