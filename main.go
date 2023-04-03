package main

import (
	"github.com/MrColorado/epubScraper/awsWrapper"
	"github.com/MrColorado/epubScraper/config"
	"github.com/MrColorado/epubScraper/scraper"
	"github.com/MrColorado/epubScraper/utils"
)

// import "github.com/MrColorado/epubScraper/cmd"

// func main() {
// 	cmd.Execute()
// }

func main() {
	config := config.GetConfig()
	client := awsWrapper.NewClient(config.AwsConfig)
	io := utils.NewS3IO(client)
	scraper := scraper.NewReadNovelScrapper(config.ScraperConfig, io)
	scraper.ScrapeNovel("pocket-hunting-dimension-v1", "")
}
