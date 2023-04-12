package main

import (
	"github.com/MrColorado/epubScraper/configuration"
	"github.com/MrColorado/epubScraper/converter"
	"github.com/MrColorado/epubScraper/dataWrapper"
	"github.com/MrColorado/epubScraper/grpcWrapper"
	"github.com/MrColorado/epubScraper/scraper"
	"github.com/MrColorado/epubScraper/utils"
)

func main() {
	config := configuration.GetConfig()
	awsClient := dataWrapper.NewAwsClient(config.AwsConfig)
	postgreClient := dataWrapper.NewPostgresClient(config.PostgresConfig)
	var io utils.IO = utils.NewS3IO(awsClient, postgreClient)

	// // Scraper
	var scraper scraper.Scraper = scraper.NewReadNovelScrapper(config.ScraperConfig, io)
	scraper.ScrapeNovel("rebirth-of-the-thief-who-roamed-the-world")

	// // Converter
	var conv converter.Converter = converter.NewEpubConverter(config.ConverterConfig, io)
	conv.ConvertNovel("rebirth-of-the-thief-who-roamed-the-world")

	server := grpcWrapper.NewSever(io)
	server.Run()
}
