package main

import (
	"context"

	"github.com/MrColorado/backend/book-handler/internal/config"
	"github.com/MrColorado/backend/book-handler/internal/converter"
	"github.com/MrColorado/backend/book-handler/internal/handler"
	"github.com/MrColorado/backend/book-handler/internal/scraper"
	"github.com/MrColorado/backend/logger"
)

var (
	scrpCfg = map[string]int{
		scraper.ReadNovelScraperName: 2,
	}
	convsName = []string{converter.EpubConverterName}
)

func main() {
	config.InitLogger()
	config := config.GetConfig()
	logger.Infof("Misc : %s", config.MiscConfig.FilesFolder)
	nats, err := handler.NewNatsClient(config.NatsConfig, context.TODO())
	if err != nil {
		logger.Info(err.Error())
		return
	}
	manager, err := handler.NewScraperManager(nats, scrpCfg, convsName)
	if err != nil {
		logger.Info(err.Error())
		return
	}
	manager.Run()
}

// func main() {
// 	config.InitLogger()
// 	config := config.GetConfig()
// 	logger.Infof("Misc : %s", config.MiscConfig.FilesFolder)

// 	db := dataStore.NewPostgresClient(config.PostgresConfig)
// 	db.GetNovelByTitle("the frozen player returns")
// }

// func main() {
// 	config.InitLogger()
// 	config := config.GetConfig()
// 	logger.Infof("Misc : %s", config.MiscConfig.FilesFolder)

// 	scrp, err := scraper.ScraperCreator(scraper.ReadNovelScraperName)
// 	if err != nil {
// 		logger.Fatalf("Failed to init scraper %s", scraper.ReadNovelScraperName)
// 	}

// 	scrp.ScrapeNovel("THE FROZEN PLAYER RETURNS")
// }
