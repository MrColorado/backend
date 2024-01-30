package converter

import (
	"fmt"

	"github.com/MrColorado/backend/bookHandler/internal/config"
	"github.com/MrColorado/backend/bookHandler/internal/core"
	"github.com/MrColorado/backend/bookHandler/internal/dataStore"
)

// Converter interface that convert novelData to e-tablet format
type Converter interface {
	ConvertNovel(novelName string) error
	ConvertPartialNovel(novelName string, startChapter int, endChapter int) error
}

func ConverterCreator(name string) (Converter, error) {
	config := config.GetConfig()
	app := core.NewApp(
		dataStore.NewAwsClient(config.AwsConfig),
		dataStore.NewPostgresClient(config.PostgresConfig),
	)

	switch name {
	case EpubConverterName:
		return NewEpubConverter(app), nil
	default:
		return nil, fmt.Errorf("failed to create converted named : %s", name)
	}
}
