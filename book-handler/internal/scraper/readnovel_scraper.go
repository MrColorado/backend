package scraper

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/MrColorado/backend/book-handler/internal/core"
	"github.com/MrColorado/backend/book-handler/internal/models"
	"github.com/MrColorado/backend/internal/common"
	"github.com/MrColorado/backend/logger"
	"github.com/gocolly/colly/v2"
)

const (
	ReadNovelScraperName = "ReadNovel"
	readNovelURL         = "https://readnovelfull.com"
)

type ReadNovelScraper struct {
	collector *colly.Collector
	app       *core.App
}

func (scraper ReadNovelScraper) findNovelURL(novelName string) (string, error) {
	url := fmt.Sprintf("%s/novel-list/search?keyword=%s", readNovelURL, strings.ReplaceAll(novelName, " ", "+"))
	nbFound := 0
	novelURL := ""

	scraper.collector.OnHTML(".list-novel", func(e *colly.HTMLElement) {
		e.ForEach(".novel-title", func(_ int, title *colly.HTMLElement) {
			nbFound += 1
			novelURL = fmt.Sprintf("%s%s", readNovelURL, title.ChildAttr("a", "href"))
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

func (scraper ReadNovelScraper) getNbOfChapter(novelID string) int {
	counter := 0

	scraper.collector.OnHTML(".panel-body", func(e *colly.HTMLElement) {
		e.ForEach("li", func(_ int, paragraph *colly.HTMLElement) {
			counter += 1
		})
	})

	defer func() {
		scraper.collector.OnHTMLDetach(".panel-body")
	}()

	scraper.collector.Visit("https://readnovelfull.com/ajax/chapter-archive?novelId=" + novelID)

	return counter
}

func (scraper ReadNovelScraper) scrapMetaData(url string, novelMetaData *models.NovelMetaData) {
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
		novelMetaData.FirstURL = readNovelURL + e.Attr("href")
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
		novelMetaData.CoverPath, _ = e.DOM.Find("img").Attr("src")
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

func (scraper ReadNovelScraper) scrapPage(url string, chapterData *models.NovelChapterData) string {
	logger.Infof("Scrape : %s", url)
	nextURL := ""

	scraper.collector.OnHTML("#next_chap", func(e *colly.HTMLElement) {
		if e.Attr("href") == "" {
			return
		}
		nextURL = readNovelURL + e.Attr("href")
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
	return nextURL
}

func NewReadNovelScrapper(app *core.App) ReadNovelScraper {
	return ReadNovelScraper{
		collector: colly.NewCollector(colly.AllowURLRevisit()),
		app:       app,
	}
}

// ScrapeNovelStart get chapter of a specific novel starting a defined chapter
func (scraper ReadNovelScraper) scrapeNovelStart(novelName string, startChapter int) {
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

	for ; url != ""; i++ {
		chapterData := models.NovelChapterData{
			Chapter: i,
		}
		url = scraper.scrapPage(url, &chapterData)

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
func (scraper ReadNovelScraper) GetName() string {
	return ReadNovelScraperName
}

// ScrapeNovel get each chapter of a specific novel
func (scraper ReadNovelScraper) ScrapeNovel(novelName string) {
	novelName = common.HarmonizeTitle(novelName)
	data, _ := scraper.app.GetMetaData(novelName)
	if data.Title == "" {
		scraper.scrapeNovelStart(novelName, 1)
	} else {
		scraper.scrapeNovelStart(novelName, data.CurrentChapter)
	}
}

// CanScrapeNovel check if novel is on the webSite
func (scraper ReadNovelScraper) CanScrapeNovel(novelName string) bool {
	novelName = common.HarmonizeTitle(novelName)
	logger.Infof("CanScrapeNovel : %s", novelName)
	novelName = strings.TrimSpace(strings.ToLower(novelName))
	novelURL, err := scraper.findNovelURL(novelName)

	if err != nil {
		return false
	}

	return len(novelURL) > 0
}
