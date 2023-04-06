package main

import (
	"github.com/MrColorado/epubScraper/awsWrapper"
	"github.com/MrColorado/epubScraper/configuration"
	"github.com/MrColorado/epubScraper/converter"
	"github.com/MrColorado/epubScraper/scraper"
	"github.com/MrColorado/epubScraper/utils"
)

func main() {
	config := configuration.GetConfig()
	client := awsWrapper.NewClient(config.AwsConfig)
	var io utils.IO = utils.NewS3IO(client)
	// var io utils.IO = utils.NewDiskIO("volumes/disk")

	// Scraper
	var scraper scraper.Scraper = scraper.NewReadNovelScrapper(config.ScraperConfig, io)
	scraper.ScrapeNovel("rebirth-of-the-thief-who-roamed-the-world")

	// Converter
	var conv converter.Converter = converter.NewEpubConverter(config.ConverterConfig, io)
	conv.ConvertNovel("rebirth-of-the-thief-who-roamed-the-world")
}
