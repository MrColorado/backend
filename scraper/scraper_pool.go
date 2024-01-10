package scraper

import (
	"fmt"
	"sync"
)

const (
	queueSize  = 5
	nbOfWorker = 2
)

type ScraperManager struct {
	mu sync.Mutex

	metaScrapers []Scraper
	scraperPools map[string]workerPool
}

func NewScraperManager() *ScraperManager {
	return &ScraperManager{
		mu:           sync.Mutex{},
		metaScrapers: []Scraper{},
		scraperPools: map[string]workerPool{},
	}
}

func (sm *ScraperManager) AddScraper(name string) error {
	fmt.Println("AddScraper")
	_, ok := sm.scraperPools[name]
	if ok {
		return fmt.Errorf("scraper %s is already register", name)
	}

	sm.mu.Lock()
	scraper, err := ScraperCreator(name)
	if err != nil {
		return fmt.Errorf("failed create pool of scraper %s", name)
	}
	sm.metaScrapers = append(sm.metaScrapers, scraper)
	sm.mu.Unlock()

	sm.scraperPools[name] = newWorkerPool(name, nbOfWorker, queueSize)
	return nil
}

func (sm *ScraperManager) CanScrape(novelName string) bool {
	canScrape := false

	sm.mu.Lock()
	for _, scraper := range sm.metaScrapers {
		if scraper.CanScrapeNovel(novelName) {
			canScrape = true
			break
		}
	}
	sm.mu.Unlock()

	return canScrape
}

func (sm *ScraperManager) Scrape(scraperName string, novelName string) error {
	wp, ok := sm.scraperPools[scraperName]
	if !ok {
		return fmt.Errorf("scraper %s does not exist", scraperName)
	}

	wp.execute(job{novelName: novelName})
	return nil
}

func (sm *ScraperManager) ShutDown() {
	for _, sp := range sm.scraperPools {
		sp.close()
	}
}
