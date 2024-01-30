package handler

import (
	"encoding/json"
	"fmt"
	"sync"

	"github.com/MrColorado/backend/bookHandler/internal/scraper"
	msgType "github.com/MrColorado/backend/internal/message"
)

type ScraperManager struct {
	nats *NatsClient

	mu           sync.Mutex
	metaScrapers []scraper.Scraper
	scraperPools map[string]WorkerPool
}

func NewScraperManager(nats *NatsClient, scraperCfg map[string]int, convsName []string) (*ScraperManager, error) {
	meta := []scraper.Scraper{}
	for name, _ := range scraperCfg {
		scrp, err := scraper.ScraperCreator(name)
		if err != nil {
			fmt.Println(err.Error())
			return nil, fmt.Errorf("failed to create scraper %s", name)
		}
		meta = append(meta, scrp)
	}

	pools := map[string]WorkerPool{}
	for scrpName, nbWorker := range scraperCfg {
		pools[scrpName] = newWorkerPool(scrpName, convsName, nbWorker, nbWorker*2)
	}

	return &ScraperManager{
		nats:         nats,
		mu:           sync.Mutex{},
		metaScrapers: meta,
		scraperPools: pools,
	}, nil
}

func (sm *ScraperManager) Run() {
	sm.nats.AddChanQueueSub("scrapable", "bookHandlerGroup")
	for scrpName, _ := range sm.scraperPools {
		sm.nats.AddChanQueueSub(fmt.Sprintf("scrape:%s", scrpName), "bookHandlerGroup")
	}

	sm.nats.Run(sm.msgHandler, sm.requestHandler)
}

func (sm *ScraperManager) msgHandler(data []byte, subject string) {
	var msg msgType.Message
	err := json.Unmarshal(data, &msg)
	if err != nil {
		fmt.Printf("failed to Unmarshal %s\n", data)
		return
	}

	switch msg.Event {
	case "scrap":
		ok, err := sm.scrape(msg.Payload)
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		if !ok {
			sm.nats.RemoveChanQueueSub(subject)
		}
	}
}

func (sm *ScraperManager) requestHandler(data []byte, _ string) ([]byte, error) {
	var msg msgType.Message
	err := json.Unmarshal(data, &msg)
	if err != nil {
		fmt.Println(err.Error())
		return generateError(1, "TODO"), fmt.Errorf("failed to Unmarshal %s\n", data)
	}

	switch msg.Event {
	case "can_scrape":
		return sm.canScrape(msg.Payload)
	}

	return generateError(3, "TODO"), nil
}

func (sm *ScraperManager) canScrape(data any) ([]byte, error) {
	rqt, ok := data.(msgType.CanScrapeRqt)
	if !ok {
		fmt.Println("failed to cast data to type msgType.CanScrapeRqt")

	}

	name := ""
	sm.mu.Lock()
	for _, scraper := range sm.metaScrapers {
		if scraper.CanScrapeNovel(rqt.Title) {
			name = scraper.GetName()
			break
		}
	}
	sm.mu.Unlock()

	msg := msgType.Message{
		Event: "can_scrape",
		Payload: msgType.CanScrapeRsp{
			ScraperName: name,
		},
	}
	rsp, err := json.Marshal(msg)
	if err != nil {
		fmt.Println(err.Error())
		return generateError(2, "TODO"), fmt.Errorf("failed to marshal CanScrapeRsp with name %s", name)
	}
	return rsp, nil
}

func (sm *ScraperManager) scrape(data any) (bool, error) {
	rqt, ok := data.(msgType.ScrapeNovelRqt)
	if !ok {
		return true, fmt.Errorf("failed to cast data to type msgType.ScrapeNovelRqt")
	}

	wp, ok := sm.scraperPools[rqt.ScraperName]
	if !ok {

		return ok, fmt.Errorf("scraper %s does not exist", rqt.ScraperName)
	}

	return wp.Execute(job{NovelName: rqt.NovelTitle}), nil
}

func generateError(code int, value string) []byte {
	msg := msgType.Message{
		Event: "error",
		Payload: msgType.Error{
			Code:  code,
			Value: value,
		},
	}

	data, err := json.Marshal(msg)
	if err != nil {
		fmt.Printf("failed to marsahl error with code : %d and value : %s", code, value)
	}
	return data
}
