package scraper

import (
	"fmt"

	"github.com/gocolly/colly"
)

const (
	vipURL string = "https://vipnovel.com/vipnovel"
)

// VipNovelScraper scrapper use to scrap https://vipnovel.com/ website
type VipNovelScraper struct{}

func (scraper *VipNovelScraper) getNextPageURL(c *colly.Collector, nextPageURL *string) {

}

// ScrapPage get the content of specific novel chapter
func (scraper *VipNovelScraper) scrapPage(c *colly.Collector, url string) (NovelData, string) {
	novelData := NovelData{}
	nextPageURL := ""

	c.OnRequest(func(r *colly.Request) {
		// fmt.Println("Visiting", r.URL)
	})

	c.OnHTML(".next_page", func(e *colly.HTMLElement) {
		nextPageURL = e.Attr("href")
	})

	c.OnHTML("h3", func(e *colly.HTMLElement) {
		novelData.Title = e.Text
	})

	c.OnHTML(".text-left", func(e *colly.HTMLElement) {
		e.ForEach("p", func(_ int, paragraph *colly.HTMLElement) {
			novelData.Paragraph = append(novelData.Paragraph, paragraph.Text)
		})
	})

	c.Visit(url)

	return novelData, nextPageURL
}

// ScrapeNovel get each chater of a specific novel
func (scraper *VipNovelScraper) ScrapeNovel(c *colly.Collector, novelName string) {
	novels := []NovelData{}

	for novelName != "" {
		fmt.Println(novelName)
		novel, nextPageURL := scraper.scrapPage(c, "https://vipnovel.com/vipnovel/the-legendary-mechanic/chapter-1113")
		novels = append(novels, novel)
		novelName = nextPageURL
	}
}
