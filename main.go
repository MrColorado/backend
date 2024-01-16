package main

import (
	"github.com/MrColorado/epubScraper/configuration"
	"github.com/MrColorado/epubScraper/converter"
	"github.com/MrColorado/epubScraper/dataWrapper"
	"github.com/MrColorado/epubScraper/utils"
)

// func main() {
// 	config := configuration.GetConfig()
// 	awsClient := dataWrapper.NewAwsClient(config.AwsConfig)
// 	postgresClient := dataWrapper.NewPostgresClient(config.PostgresConfig)

// 	io := utils.NewS3IO(awsClient, postgresClient)
// 	var conv converter.Converter = converter.NewEpubConverter(config.ConverterConfig, io)

// 	sm := scraper.NewScraperManager()
// 	err := sm.AddScraper(scraper.ReadNovelScraperName)
// 	if err != nil {
// 		fmt.Println(err.Error())
// 		return
// 	}

// 	server := grpcWrapper.NewSever(io, sm, conv)

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

// func main() {
// 	config := configuration.GetConfig()
// 	postgresClient := dataWrapper.NewPostgresClient(config.PostgresConfig)
// 	postgresClient.GetNovelById("841b35a6-5059-41a8-a18e-8e33fa990893")
// }

func main() {
	config := configuration.GetConfig()
	awsClient := dataWrapper.NewAwsClient(config.AwsConfig)
	postgresClient := dataWrapper.NewPostgresClient(config.PostgresConfig)

	io := utils.NewS3IO(awsClient, postgresClient)
	var conv converter.Converter = converter.NewEpubConverter(config.ConverterConfig, io)

	conv.ConvertNovel("big life")
}

// "/s3/novels/under the oak tree/epub/under the oak tree-0001-0100.epub"

// func main() {
// config := configuration.GetConfig()
// awsClient := dataWrapper.NewAwsClient(config.AwsConfig)
// postgresClient := dataWrapper.NewPostgresClient(config.PostgresConfig)

// io := utils.NewS3IO(awsClient, postgresClient)
// var conv converter.Converter = converter.NewEpubConverter(config.ConverterConfig, io)
// var scraper scraper.Scraper = scraper.NewReadNovelScrapper(config.ScraperConfig, io)

// novelName := "rebirth of an idle noblewoman"
// tmp, _ := postgresClient.GetNovelByTitle(novelName)
// fmt.Printf(tmp.CoreData.Author)

// sm := scraper.NewScraperManager()
// err := sm.AddScraper(scraper.ReadNovelScraperName)
// if err != nil {
// 	fmt.Println(err.Error())
// 	return
// }
// 	sm.Scrape(scraper.ReadNovelScraperName, "inadvertently invincible")
// 	sm.Scrape(scraper.ReadNovelScraperName, "hot farmer's wife buying a husband for the farm")

// 	time.Sleep(time.Minute)

// 	sm.ShutDown()
// }
