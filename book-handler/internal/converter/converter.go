package converter

import (
	"github.com/MrColorado/backend/book-handler/internal/config"
	"github.com/MrColorado/backend/book-handler/internal/core"
	"github.com/MrColorado/backend/book-handler/internal/data"
	"github.com/MrColorado/backend/logger"
)

// Converter interface that convert novelData to e-tablet format
type Converter interface {
	ConvertNovel(novelName string) error
	ConvertPartialNovel(novelName string, startChapter int, endChapter int) error
}

func ConverterCreator(name string) (Converter, error) {
	cfg := config.GetConfig()
	app := core.NewApp(
		data.NewAwsClient(cfg.AwsConfig),
		data.NewPostgresClient(cfg.PostgresConfig),
	)

	switch name {
	case EpubConverterName:
		return NewEpubConverter(app), nil
	default:
		return nil, logger.Errorf("failed to create converted named : %s", name)
	}
}
