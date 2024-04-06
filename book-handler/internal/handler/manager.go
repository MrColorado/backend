package handler

import (
	"encoding/json"
	"fmt"
	"sync"

	"github.com/MrColorado/backend/book-handler/internal/scraper"
	msgType "github.com/MrColorado/backend/internal/message"
	"github.com/MrColorado/backend/logger"
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
			return nil, logger.Errorf("failed to create scraper %s", name)
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
		sm.nats.AddChanQueueSub(fmt.Sprintf("scraper.%s", scrpName), "bookHandlerGroup")
	}

	sm.nats.Run(sm.msgHandler, sm.requestHandler)
}

func (sm *ScraperManager) msgHandler(data []byte, subject string) {
	logger.Infof("Data : %s | From : %s", string(data), subject)
	var msg msgType.Message
	err := json.Unmarshal(data, &msg)
	if err != nil {
		logger.Warnf("failed to Unmarshal %s : %s", data, err.Error())
		return
	}

	switch msg.Event {
	case "scrape":
		ok, err := sm.scrape(msg.Payload)
		if err != nil {
			return
		}

		if !ok {
			sm.nats.RemoveChanQueueSub(subject) // TODO how de we subscribe again after ?
		}
	}
}

func (sm *ScraperManager) requestHandler(data []byte, subject string) ([]byte, error) {
	logger.Infof("Data : %s | From : %s", string(data), subject)
	var msg msgType.Message
	err := json.Unmarshal(data, &msg)
	if err != nil {
		return generateError(1, "TODO"), logger.Errorf("failed to Unmarshal %s", data)
	}

	switch msg.Event {
	case "can_scrape":
		return sm.canScrape(msg.Payload)
	}

	return generateError(3, "TODO"), nil
}

func (sm *ScraperManager) canScrape(data json.RawMessage) ([]byte, error) {
	var rqt msgType.CanScrapeRqt
	err := json.Unmarshal(data, &rqt)
	if err != nil {
		return nil, logger.Error("failed to cast data to type msgType.CanScrapeRqt")
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

	j, err := json.Marshal(msgType.CanScrapeRsp{
		ScraperName: name,
	})
	if err != nil {
		return nil, logger.Error("failed to marshal msgType.CanScrapeRsp")
	}

	msg := msgType.Message{
		Event:   "can_scrape",
		Payload: json.RawMessage(j),
	}

	rsp, err := json.Marshal(msg)
	if err != nil {
		return generateError(2, "TODO"), logger.Errorf("failed to marshal Message")
	}
	return rsp, nil
}

func (sm *ScraperManager) scrape(data json.RawMessage) (bool, error) {
	var rqt msgType.ScrapeNovelRqt
	err := json.Unmarshal(data, &rqt)
	if err != nil {
		return true, logger.Error("failed to cast data to type msgType.ScrapeNovelRqt")
	}

	wp, ok := sm.scraperPools[rqt.ScraperName]
	if !ok {
		return true, logger.Errorf("scraper %s does not exist", rqt.ScraperName)
	}

	return wp.Execute(job{NovelName: rqt.NovelTitle}), nil
}

func generateError(code int, value string) []byte {
	j, err := json.Marshal(msgType.Error{
		Code:  code,
		Value: value,
	})
	if err != nil {
		logger.Warnf("failed to marshal error with code : %d and value : %s", code, value)
	}

	msg := msgType.Message{
		Event:   "error",
		Payload: json.RawMessage(j),
	}

	data, err := json.Marshal(msg)
	if err != nil {
		logger.Warnf("failed to marshal message")
	}
	return data
}
