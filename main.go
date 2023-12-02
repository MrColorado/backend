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
	var scraper scraper.Scraper = scraper.NewReadNovelScrapper(config.ScraperConfig, io)
	var conv converter.Converter = converter.NewEpubConverter(config.ConverterConfig, io)

	server := grpcWrapper.NewSever(io, scraper, conv)
	server.Run()
}

// func main() {
// 	config := configuration.GetConfig()
// 	awsClient := dataWrapper.NewAwsClient(config.AwsConfig)
// 	postgreClient := dataWrapper.NewPostgresClient(config.PostgresConfig)

// 	var io utils.IO = utils.NewS3IO(awsClient, postgreClient)
// 	var scraper scraper.Scraper = scraper.NewReadNovelScrapper(config.ScraperConfig, io)
// 	var conv converter.Converter = converter.NewEpubConverter(config.ConverterConfig, io)

// 	scraper.ScrapeNovel("monster paradise")
// 	conv.ConvertNovel("monster paradise")
// }
