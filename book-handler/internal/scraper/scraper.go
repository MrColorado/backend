package scraper

import (
	"fmt"

	"github.com/MrColorado/backend/book-handler/internal/config"
	"github.com/MrColorado/backend/book-handler/internal/core"
	"github.com/MrColorado/backend/book-handler/internal/data"
)

// Scraper of each website should implement this interface
type Scraper interface {
	GetName() string
	ScrapeNovel(novelName string)
	CanScrapeNovel(novelName string) bool
}

func ScraperCreator(scraperName string) (Scraper, error) {
	cfg := config.GetConfig()
	app := core.NewApp(
		data.NewAwsClient(cfg.AwsConfig),
		data.NewPostgresClient(cfg.PostgresConfig),
	)

	switch scraperName {
	case ReadNovelScraperName:
		return NewReadNovelScrapper(app), nil
	default:
		return nil, fmt.Errorf("failed to create scraper named : %s", scraperName)
	}
}
