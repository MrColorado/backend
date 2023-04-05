package main

import (
	"github.com/MrColorado/epubScraper/awsWrapper"
	"github.com/MrColorado/epubScraper/configuration"
	"github.com/MrColorado/epubScraper/scraper"
	"github.com/MrColorado/epubScraper/utils"
)

func main() {
	config := configuration.GetConfig()
	client := awsWrapper.NewClient(config.AwsConfig)
	var io utils.IO = utils.NewS3IO(client)
	var scraper scraper.Scraper = scraper.NewReadNovelScrapper(config.ScraperConfig, io)
	scraper.ScrapeNovel("the-frozen-player-returns")
}
