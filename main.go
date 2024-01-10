package main

import (
	"fmt"
	"time"

	"github.com/MrColorado/epubScraper/scraper"
)

// func main() {
// 	config := configuration.GetConfig()
// 	awsClient := dataWrapper.NewAwsClient(config.AwsConfig)
// 	postgresClient := dataWrapper.NewPostgresClient(config.PostgresConfig)

// 	io := utils.NewS3IO(awsClient, postgresClient)
// 	var scraper scraper.Scraper = scraper.NewReadNovelScrapper(config.ScraperConfig, io)
// 	var conv converter.Converter = converter.NewEpubConverter(config.ConverterConfig, io)

// 	server := grpcWrapper.NewSever(io, scraper, conv)
// 	server.Run()
// }

// func main() {
// 	config := configuration.GetConfig()
// 	awsClient := dataWrapper.NewAwsClient(config.AwsConfig)
// 	postgresClient := dataWrapper.NewPostgresClient(config.PostgresConfig)

// 	io := utils.NewS3IO(awsClient, postgresClient)
// 	io.GetBook("under the oak tree", 1, 100)

// 	// awsClient.GetPreSignedLink("gods' impact online/cover.jpg")
// }

// "/s3/novels/under the oak tree/epub/under the oak tree-0001-0100.epub"

func main() {
	// config := configuration.GetConfig()
	// awsClient := dataWrapper.NewAwsClient(config.AwsConfig)
	// postgresClient := dataWrapper.NewPostgresClient(config.PostgresConfig)

	// io := utils.NewS3IO(awsClient, postgresClient)
	// var conv converter.Converter = converter.NewEpubConverter(config.ConverterConfig, io)
	// var scraper scraper.Scraper = scraper.NewReadNovelScrapper(config.ScraperConfig, io)

	// novelName := "rebirth of an idle noblewoman"
	// tmp, _ := postgresClient.GetNovelByTitle(novelName)
	// fmt.Printf(tmp.CoreData.Author)

	sm := scraper.NewScraperManager()
	err := sm.AddScraper(scraper.ReadNovelScraperName)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	// if sm.CanScrape("inadvertently invincible") {
	// 	fmt.Println("ok")
	// }
	// if sm.CanScrape("hot farmer's wife buying a husband for the farm") {
	// 	fmt.Println("ok")
	// }

	sm.Scrape(scraper.ReadNovelScraperName, "inadvertently invincible")
	sm.Scrape(scraper.ReadNovelScraperName, "hot farmer's wife buying a husband for the farm")

	time.Sleep(time.Minute)

	sm.ShutDown()
}
