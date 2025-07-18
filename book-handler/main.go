package main

import (
	"github.com/MrColorado/backend/book-handler/internal/config"
	"github.com/MrColorado/backend/book-handler/internal/converter"
	"github.com/MrColorado/backend/book-handler/internal/scraper"
	"github.com/MrColorado/backend/logger"
)

var (
	scrpCfg = map[string]int{
		scraper.NovelBinScraperName:  2,
		scraper.ReadNovelScraperName: 2,
	}
	convsName = []string{converter.EpubConverterName}
)

// func main() {
// 	config.InitLogger()
// 	cfg := config.GetConfig()
// 	logger.Infof("Misc : %s", cfg.MiscConfig.FilesFolder)
// 	nats, err := handler.NewNatsClient(cfg.NatsConfig, context.TODO())
// 	if err != nil {
// 		logger.Info(err.Error())
// 		return
// 	}
// 	manager, err := handler.NewScraperManager(nats, scrpCfg, convsName)
// 	if err != nil {
// 		logger.Info(err.Error())
// 		return
// 	}
// 	manager.Run()
// }

// func main() {
// 	config.InitLogger()
// 	config := config.GetConfig()
// 	logger.Infof("Misc : %s", config.MiscConfig.FilesFolder)

// 	db := data.NewPostgresClient(config.PostgresConfig)
// 	db.GetBookByTitle("big life", 101)
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

func main() {
	config.InitLogger()
	cfg := config.GetConfig()
	logger.Infof("Misc : %s", cfg.MiscConfig.FilesFolder)

	scrp, err := scraper.ScraperCreator(scraper.NovelBinScraperName)
	if err != nil {
		logger.Fatalf("Failed to init scraper %s", scraper.NovelBinScraperName)
	}

	scrp.ScrapeNovel("defiance of the fall")
}

// func main() {
// 	config.InitLogger()
// 	cfg := config.GetConfig()
// 	logger.Infof("Misc : %s", cfg.MiscConfig.FilesFolder)

// 	db := data.NewPostgresClient(cfg.PostgresConfig)
// 	db.GetBookByTitle("big life", 101)

// 	conv, err := converter.ConverterCreator(converter.EpubConverterName)
// 	if err != nil {
// 		logger.Fatal(err.Error())
// 	}

// 	conv.ConvertNovel("genetic ascension")
// }
