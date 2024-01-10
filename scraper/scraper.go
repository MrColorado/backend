package scraper

import (
	"fmt"

	"github.com/MrColorado/epubScraper/configuration"
	"github.com/MrColorado/epubScraper/dataWrapper"
	"github.com/MrColorado/epubScraper/utils"
)

// Scraper of each website should implement this interface
type Scraper interface {
	ScrapeNovel(novelName string)
	CanScrapeNovel(novelName string) bool
}

func ScraperCreator(scraperName string) (Scraper, error) {
	config := configuration.GetConfig()
	io := utils.NewS3IO(
		dataWrapper.NewAwsClient(config.AwsConfig),
		dataWrapper.NewPostgresClient(config.PostgresConfig),
	)

	switch scraperName {
	case ReadNovelScraperName:
		return NewReadNovelScrapper(config.ScraperConfig, io), nil
	default:
		return nil, fmt.Errorf("failed to create scraper named : %s", scraperName)
	}
}
