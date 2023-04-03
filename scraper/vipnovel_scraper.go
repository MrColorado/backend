package scraper

// import (
// 	"fmt"

// 	"github.com/MrColorado/epubScraper/utils"
// 	"github.com/gocolly/colly"
// )

// const (
// 	vipURL string = "https://vipnovel.com/vipnovel"
// )

// // VipNovelScraper scrapper use to scrap https://vipnovel.com/ website
// type VipNovelScraper struct{}

// func (scraper *VipNovelScraper) getNextPageURL(c *colly.Collector, nextPageURL *string) {

// }

// // ScrapPage get the content of specific novel chapter
// func (scraper *VipNovelScraper) scrapPage(c *colly.Collector, url string) (utils.NovelChapterData, string) {
// 	novelData := utils.NovelChapterData{}
// 	nextPageURL := ""

// 	c.OnHTML(".next_page", func(e *colly.HTMLElement) {
// 		nextPageURL = e.Attr("href")
// 	})

// 	c.OnHTML(".text-left", func(e *colly.HTMLElement) {
// 		e.ForEach("p", func(_ int, paragraph *colly.HTMLElement) {
// 			novelData.Paragraph = append(novelData.Paragraph, paragraph.Text)
// 		})
// 	})

// 	c.Visit(url)

// 	return novelData, nextPageURL
// }

// // ScrapeNovel get each chater of a specific novel
// func (scraper *VipNovelScraper) ScrapeNovel(c *colly.Collector, novelName string) {
// 	novels := []utils.NovelChapterData{}

// 	for novelName != "" {
// 		fmt.Println(novelName)
// 		novel, nextPageURL := scraper.scrapPage(c, "https://vipnovel.com/vipnovel/the-legendary-mechanic/chapter-1113")
// 		novels = append(novels, novel)
// 		novelName = nextPageURL
// 	}
// }
