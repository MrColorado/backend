package scraper

import (
	"fmt"

	"github.com/MrColorado/backend/bookHandler/internal/config"
	"github.com/MrColorado/backend/bookHandler/internal/core"
	"github.com/MrColorado/backend/bookHandler/internal/dataStore"
)

// Scraper of each website should implement this interface
type Scraper interface {
	GetName() string
	ScrapeNovel(novelName string)
	CanScrapeNovel(novelName string) bool
}

func ScraperCreator(scraperName string) (Scraper, error) {
	config := config.GetConfig()
	app := core.NewApp(
		dataStore.NewAwsClient(config.AwsConfig),
		dataStore.NewPostgresClient(config.PostgresConfig),
	)

	switch scraperName {
	case ReadNovelScraperName:
		return NewReadNovelScrapper(app), nil
	default:
		return nil, fmt.Errorf("failed to create scraper named : %s", scraperName)
	}
}
