package main

import (
	"github.com/MrColorado/epubScraper/converter"
)

func main() {
	// scraper := scraper.ReadNovelScraper{}
	// c := colly.NewCollector()
	// scraper.ScrapPartialNovel(c, "warlock-of-the-magus-world", 684, 684)
	// scraper.ScrapeNovel(c, "rebuild-world")

	converter := converter.EpubConverter{}
	// converter.ConvertPartialNovel("/home/mrcolorado/Novels/raw/warlock-of-the-magus-world",
	// 	"/home/mrcolorado/Novels/epub/warlock-of-the-magus-world",
	// 	684, 784)
	converter.ConvertPartialNovel("/home/mrcolorado/Novels/raw/warlock-of-the-magus-world",
		"/home/mrcolorado/Novels/epub/warlock-of-the-magus-world", 700, 850)

	// for _, novelData := range novelsData {
	// 	fmt.Println("---------------------------------------------------")
	// 	fmt.Println(novelData.Title)
	// 	for _, paragraph := range novelData.Paragraph {
	// 		fmt.Println(paragraph)
	// 	}
	// 	fmt.Println("---------------------------------------------------")
	// }
}
