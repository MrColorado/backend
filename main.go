package main

import "github.com/MrColorado/epubScraper/converter"

func main() {
	// scraper := scraper.ReadNovelScraper{}
	// c := colly.NewCollector()
	// scraper.ScrapPartialNovel(c, "warlock-of-the-magus-world", 1, 1, "/home/mrcolorado/Novels/raw")
	//scraper.ScrapeNovel(c, "warlock-of-the-magus-world", "/home/mrcolorado/Novels/raw")

	converter := converter.EpubConverter{}
	// converter.ConvertPartialNovel("/home/mrcolorado/Novels/raw/warlock-of-the-magus-world",
	// 	"/home/mrcolorado/Novels/epub/warlock-of-the-magus-world",
	// 	1, 50)
	converter.ConvertNovel("/home/mrcolorado/Novels/raw/warlock-of-the-magus-world",
		"/home/mrcolorado/Novels/epub/warlock-of-the-magus-world")

	// for _, novelData := range novelsData {
	// 	fmt.Println("---------------------------------------------------")
	// 	fmt.Println(novelData.Title)
	// 	for _, paragraph := range novelData.Paragraph {
	// 		fmt.Println(paragraph)
	// 	}
	// 	fmt.Println("---------------------------------------------------")
	// }
}
