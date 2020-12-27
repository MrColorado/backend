package main

import (
	"github.com/MrColorado/epubScraper/scraper"
	"github.com/gocolly/colly"
)

func main() {
	scraper := scraper.ReadNovelScraper{}

	c := colly.NewCollector()
	// scraper.ScrapPartialNovel(c, "warlock-of-the-magus-world", 684, 1200)
	scraper.ScrapeNovel(c, "rebuild-world")

	// for _, novelData := range novelsData {
	// 	fmt.Println("---------------------------------------------------")
	// 	fmt.Println(novelData.Title)
	// 	for _, paragraph := range novelData.Paragraph {
	// 		fmt.Println(paragraph)
	// 	}
	// 	fmt.Println("---------------------------------------------------")
	// }
}
