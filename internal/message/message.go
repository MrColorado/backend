package msgType

type CanScrapeRqt struct {
	Title string `json:"title"`
}

type CanScrapeRsp struct {
	ScraperName string `json:"scraper_name"`
}

type ScrapeNovelRqt struct {
	ScraperName string `json:"scraper_name"`
	NovelTitle  string `json:"novel_title"`
}
