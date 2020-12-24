package scraper

import (
	"github.com/gocolly/colly"
)

type NovelData struct {
	Title     string
	Chapter   string
	Paragraph []string
}

type Scraper interface {
	ScrapeNovel(c *colly.Collector, novelName string) []NovelData
	ScrapPartialNovel(c *colly.Collector, novelName string, startChapter int, endChapter int) []NovelData
}
