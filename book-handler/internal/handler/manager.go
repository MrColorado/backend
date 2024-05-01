package handler

import (
	"encoding/json"
	"sync"
	"time"

	"github.com/MrColorado/backend/book-handler/internal/scraper"
	msgType "github.com/MrColorado/backend/internal/message"
	"github.com/MrColorado/backend/logger"
)

type ScraperManager struct {
	scraperPools map[string]WorkerPool
	nats         *NatsClient
	metaScrapers []scraper.Scraper
	mu           sync.Mutex
}

func NewScraperManager(nats *NatsClient, scraperCfg map[string]int, convsName []string) (*ScraperManager, error) {
	meta := []scraper.Scraper{}
	for name := range scraperCfg {
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
	for scrpName := range sm.scraperPools {
		sm.nats.AddChanQueueSub("scraper."+scrpName, "bookHandlerGroup")
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

	if msg.Event == "scrape" {
		var rqt msgType.ScrapeNovelRqt
		err := json.Unmarshal(msg.Payload, &rqt)
		if err != nil {
			logger.Error("failed to cast data to type msgType.ScrapeNovelRqt")
			return
		}

		canQueue, err := sm.scrape(rqt)
		if err != nil {
			return
		}

		if !canQueue {
			sm.nats.RemoveChanQueueSub(subject, false)
			sm.delaySub(subject, rqt.ScraperName)
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

	if msg.Event == "can_scrape" {
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
	for _, scrp := range sm.metaScrapers {
		if scrp.CanScrapeNovel(rqt.Title) {
			name = scrp.GetName()
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

func (sm *ScraperManager) scrape(rqt msgType.ScrapeNovelRqt) (bool, error) {
	wp, ok := sm.scraperPools[rqt.ScraperName]
	if !ok {
		return true, logger.Errorf("scraper %s does not exist", rqt.ScraperName)
	}
	wp.Execute(job{NovelName: rqt.NovelTitle})

	return wp.CanQueuJob(), nil
}

func (sm *ScraperManager) delaySub(subject string, scraperName string) {
	logger.Info("Start delaySub")
	ticker := time.NewTicker(500 * time.Millisecond)

	tempo := make(chan bool)

	go func() {
		tempo <- true
		for {
			select {
			// Add ctx to cancel go routing & add waitgroup
			case <-ticker.C:
				logger.Infof("Tick for scraper : %s", scraperName)

				wp, ok := sm.scraperPools[scraperName]
				if !ok {
					logger.Errorf("scraper %s does not exist", scraperName)
					return
				}

				if wp.CanQueuJob() {
					sm.nats.AddChanQueueSub(subject, "bookHandlerGroup")
					ticker.Stop()
					return
				}
			}
		}
	}()

	<-tempo
	close(tempo)
	logger.Info("End delaySub")
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
