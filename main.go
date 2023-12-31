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
	postgresClient := dataWrapper.NewPostgresClient(config.PostgresConfig)

	io := utils.NewS3IO(awsClient, postgresClient)
	var scraper scraper.Scraper = scraper.NewReadNovelScrapper(config.ScraperConfig, io)
	var conv converter.Converter = converter.NewEpubConverter(config.ConverterConfig, io)

	server := grpcWrapper.NewSever(io, scraper, conv)
	server.Run()
}

// func main() {
// 	config := configuration.GetConfig()
// 	awsClient := dataWrapper.NewAwsClient(config.AwsConfig)
// 	postgresClient := dataWrapper.NewPostgresClient(config.PostgresConfig)

// 	io := utils.NewS3IO(awsClient, postgresClient)
// 	io.GetBook("under the oak tree", 1, 100)
// }

// "/s3/novels/under the oak tree/epub/under the oak tree-0001-0100.epub"

// func main() {
// 	config := configuration.GetConfig()
// 	awsClient := dataWrapper.NewAwsClient(config.AwsConfig)
// 	postgresClient := dataWrapper.NewPostgresClient(config.PostgresConfig)

// 	io := utils.NewS3IO(awsClient, postgresClient)
// 	var scraper scraper.Scraper = scraper.NewReadNovelScrapper(config.ScraperConfig, io)
// 	var conv converter.Converter = converter.NewEpubConverter(config.ConverterConfig, io)

// 	novelName := "the steward demonic emperor"
// 	scraper.ScrapeNovel(novelName)
// 	conv.ConvertNovel(novelName)
// }
