package scraper

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/MrColorado/backend/book-handler/internal/core"
	"github.com/MrColorado/backend/book-handler/internal/models"
	"github.com/MrColorado/backend/internal/common"
	"github.com/MrColorado/backend/logger"
	"github.com/gocolly/colly/v2"
)

const (
	NovelBinScraperName = "NovelBin"
	novelBinURL         = "https://novelbin.me"
)

type NovelBinScraper struct {
	collector *colly.Collector
	app       *core.App
}

func (scraper NovelBinScraper) findNovelURL(novelName string) (string, error) {
	url := fmt.Sprintf("%s/search?keyword=%s", novelBinURL, strings.ReplaceAll(novelName, " ", "+"))
	nbFound := 0
	novelURL := ""

	scraper.collector.OnHTML(".list-novel", func(e *colly.HTMLElement) {
		e.ForEach(".novel-title", func(_ int, title *colly.HTMLElement) {
			nbFound += 1
			novelURL = title.ChildAttr("a", "href")
		})
	})

	err := scraper.collector.Visit(url)
	if err != nil {
		return "", logger.Errorf("failed to visit url %s : %s", novelURL, err.Error())
	}

	if nbFound == 1 {
		return novelURL, nil
	}
	return "", nil
}

func (scraper NovelBinScraper) getNbOfChapter(novelID string) int {
	counter := 0

	scraper.collector.OnHTML(".panel-body", func(e *colly.HTMLElement) {
		e.ForEach("li", func(_ int, paragraph *colly.HTMLElement) {
			counter += 1
		})
	})

	defer func() {
		scraper.collector.OnHTMLDetach(".panel-body")
	}()

	scraper.collector.Visit("https://novelbin.me/ajax/chapter-archive?novelId=" + novelID)

	return counter
}

func (scraper NovelBinScraper) scrapMetaData(url string, novelMetaData *models.NovelMetaData) {
	novelID := ""
	novelMetaData.CurrentChapter = 1

	scraper.collector.OnHTML("#tab-description", func(e *colly.HTMLElement) {
		e.ForEach("p", func(_ int, paragraph *colly.HTMLElement) {
			novelMetaData.Summary += paragraph.Text
		})
	})

	scraper.collector.OnHTML("#rating", func(e *colly.HTMLElement) {
		novelID = e.Attr("data-novel-id")
	})

	scraper.collector.OnHTML(".title", func(e *colly.HTMLElement) {
		novelMetaData.Title = common.HarmonizeTitle(e.Text)
	})

	scraper.collector.OnHTML(".btn-read-now", func(e *colly.HTMLElement) {
		novelMetaData.FirstURL = e.Attr("href")
	})

	scraper.collector.OnHTML(".info-meta", func(e *colly.HTMLElement) {
		e.ForEach("li", func(_ int, li *colly.HTMLElement) {
			if title := li.ChildText("h3"); title == "Author:" {
				li.ForEach("a", func(_ int, a *colly.HTMLElement) {
					novelMetaData.Author = a.Text // TODO get all name or only the first ?
				})
			} else if genre := li.ChildText("h3"); genre == "Genre:" {
				li.ForEach("a", func(_ int, a *colly.HTMLElement) {
					novelMetaData.Genres = append(novelMetaData.Genres, a.Text)
				})
			} else if status := li.ChildText("h3"); status == "Status:" {
				li.ForEach("a", func(_ int, a *colly.HTMLElement) {
					novelMetaData.Status = 1
					if a.Text == "Ongoing" {
						novelMetaData.Status = 0
					}
				})
			}
		})
	})

	scraper.collector.OnHTML(".book", func(e *colly.HTMLElement) {
		tmp := e.DOM.Find("img")
		novelMetaData.CoverPath, _ = tmp.Attr("data-src")
	})

	defer func() {
		scraper.collector.OnHTMLDetach("#tab-description")
		scraper.collector.OnHTMLDetach(".btn-read-now")
		scraper.collector.OnHTMLDetach(".next_chap")
		scraper.collector.OnHTMLDetach(".title")
		scraper.collector.OnHTMLDetach(".book")
	}()

	scraper.collector.Visit(url)

	novelMetaData.NbChapter = scraper.getNbOfChapter(novelID)
}

func (scraper NovelBinScraper) scrapPage(url string, chapterData *models.NovelChapterData) (string, error) {
	logger.Infof("Scrape : %s", url)
	nextURL := ""

	scraper.collector.OnHTML("#next_chap", func(e *colly.HTMLElement) {
		if e.Attr("href") == "" {
			return
		}
		nextURL = e.Attr("href")
	})

	scraper.collector.OnHTML("#chr-content", func(e *colly.HTMLElement) {
		e.ForEach("p", func(_ int, paragraph *colly.HTMLElement) {
			chapterData.Paragraph = append(chapterData.Paragraph, paragraph.Text)
		})
	})
	defer func() {
		scraper.collector.OnHTMLDetach("#chr-content")
		scraper.collector.OnHTMLDetach("#next_chap")
	}()

	scraper.collector.Visit(url)

	if nextURL == "" {
		return "", logger.Errorf("Failed to get next url, current url is : %s", url)
	}

	return nextURL, nil
}

func NewBinNovelScrapper(app *core.App) NovelBinScraper {
	logger.Infof("Create %s's scraper", NovelBinScraperName)
	return NovelBinScraper{
		collector: colly.NewCollector(colly.AllowURLRevisit()),
		app:       app,
	}
}

// ScrapeNovelStart get chapter of a specific novel starting a defined chapter
func (scraper NovelBinScraper) scrapeNovelStart(novelName string, startChapter int) {
	data, _ := scraper.app.GetMetaData(novelName)

	if data.Title == "" {
		novelURL, _ := scraper.findNovelURL(novelName)
		scraper.scrapMetaData(novelURL, &data)

		resp, err := http.Get(data.CoverPath)
		if err == nil {
			data.CoverData, _ = io.ReadAll(resp.Body)
			defer resp.Body.Close()
		}

		if scraper.app.ExportMetaData(data.Title, data, true) != nil {
			return
		}
	}

	i := 1
	url := data.FirstURL
	if startChapter != 1 {
		i = data.CurrentChapter
		url = data.NextURL
	}

	for ; i != data.NbChapter; i++ {
		chapterData := models.NovelChapterData{
			Chapter: i,
		}

		retry := 0
		for ; retry != 5; retry++ {
			nextURL, err := scraper.scrapPage(url, &chapterData)
			if err == nil {
				url = nextURL
				break
			}
			logger.Infof("Failed to get novel at URL %s for chapter %d retry in %d sec", url, i, retry)
			time.Sleep(time.Second * time.Duration(retry))
		}
		if retry == 5 {
			logger.Errorf("Failed to get novel at URL %s for chapter %d", url, i)
			return
		}

		data.NextURL = url
		data.CurrentChapter += 1
		if scraper.app.ExportNovelChapter(data.Title, chapterData) != nil {
			return
		}
		if scraper.app.ExportMetaData(data.Title, data, false) != nil {
			return
		}
	}
}

// GetName return name of scraper
func (scraper NovelBinScraper) GetName() string {
	return NovelBinScraperName
}

// ScrapeNovel get each chapter of a specific novel
func (scraper NovelBinScraper) ScrapeNovel(novelName string) {
	novelName = common.HarmonizeTitle(novelName)
	data, _ := scraper.app.GetMetaData(novelName)
	if data.Title == "" {
		scraper.scrapeNovelStart(novelName, 1)
	} else {
		scraper.scrapeNovelStart(novelName, data.CurrentChapter)
	}
}

// CanScrapeNovel check if novel is on the webSite
func (scraper NovelBinScraper) CanScrapeNovel(novelName string) bool {
	novelName = common.HarmonizeTitle(novelName)
	logger.Infof("CanScrapeNovel : %s", novelName)
	novelName = strings.TrimSpace(strings.ToLower(novelName))
	novelURL, err := scraper.findNovelURL(novelName)

	if err != nil {
		return false
	}

	return len(novelURL) > 0
}
