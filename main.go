package main

import (
	"fmt"

	"github.com/MrColorado/epubScraper/scraper"
	"github.com/gocolly/colly"
)

func main() {
	scraper := scraper.VipNovelScraper{}

	c := colly.NewCollector()
	novelsData := scraper.ScrapeNovel(c, "toto")

	for _, novelData := range novelsData {
		fmt.Println("---------------------------------------------------")
		fmt.Println(novelData.Title)
		//for _, paragraph := range novelData.Paragraph {
		//	fmt.Println(paragraph)
		//}
		fmt.Println("---------------------------------------------------")
	}
}
