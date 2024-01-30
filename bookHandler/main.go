package main

import (
	"context"
	"fmt"

	"github.com/MrColorado/backend/bookHandler/internal/config"
	"github.com/MrColorado/backend/bookHandler/internal/converter"
	"github.com/MrColorado/backend/bookHandler/internal/handler"
	"github.com/MrColorado/backend/bookHandler/internal/scraper"
)

var (
	scrpCfg = map[string]int{
		scraper.ReadNovelScraperName: 2,
	}
	convsName = []string{converter.EpubConverterName}
)

func main() {
	config := config.GetConfig()
	nats, err := handler.NewNatsClient(config.NatsConfig, context.TODO())
	if err != nil {
		fmt.Printf(err.Error())
		return
	}
	manager, err := handler.NewScraperManager(nats, scrpCfg, convsName)
	if err != nil {
		fmt.Printf(err.Error())
		return
	}
	manager.Run()
}
