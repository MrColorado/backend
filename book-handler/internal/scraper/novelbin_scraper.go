package scraper

import (
	"fmt"
	"io"
	"math/rand"
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

	err := scraper.collector.Visit("https://novelbin.me/ajax/chapter-archive?novelId=" + novelID)
	if err != nil {
		logger.Errorf("Failed to get chapter count for novel ID %s: %v", novelID, err)
		return 0
	}

	return counter
}

func (scraper NovelBinScraper) scrapMetaData(url string, novelMetaData *models.NovelMetaData) {
	logger.Infof("Scraping metadata from URL: %s", url)
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
	logger.Infof("Scraping chapter page: %s", url)
	nextURL := ""

	// Add response status checking
	scraper.collector.OnResponse(func(r *colly.Response) {
		logger.Infof("Response status: %d for URL: %s", r.StatusCode, r.Request.URL)
		if r.StatusCode != http.StatusOK {
			logger.Errorf("HTTP %d error for URL: %s", r.StatusCode, r.Request.URL)
		}
	})

	// Add error handling
	scraper.collector.OnError(func(r *colly.Response, err error) {
		logger.Errorf("Request failed for URL %s: %v", r.Request.URL, err)
	})

	// Add diagnostic logging
	scraper.collector.OnHTML("body", func(e *colly.HTMLElement) {
		logger.Infof("Chapter page body length: %d", len(e.Text))
		if len(e.Text) > 500 {
			logger.Infof("Chapter page HTML preview: %s", e.Text[:500])
		} else {
			logger.Infof("Chapter page HTML preview: %s", e.Text)
		}

		// Check if we got a 404 or error page
		if strings.Contains(e.Text, "404") || strings.Contains(e.Text, "Not Found") {
			logger.Errorf("Page not found (404) for URL: %s", url)
		}
		if strings.Contains(e.Text, "Access Denied") || strings.Contains(e.Text, "Forbidden") {
			logger.Errorf("Access denied (403) for URL: %s", url)
		}
	})

	scraper.collector.OnHTML("#next_chap", func(e *colly.HTMLElement) {
		logger.Infof("Found #next_chap element")
		if e.Attr("href") == "" {
			logger.Infof("Next chapter href is empty")
			return
		}
		nextURL = e.Attr("href")
		logger.Infof("Next URL: %s", nextURL)
	})

	scraper.collector.OnHTML("#chr-content", func(e *colly.HTMLElement) {
		logger.Infof("Found #chr-content element")
		e.ForEach("p", func(_ int, paragraph *colly.HTMLElement) {
			chapterData.Paragraph = append(chapterData.Paragraph, paragraph.Text)
		})
		logger.Infof("Found %d paragraphs in chapter", len(chapterData.Paragraph))
	})

	// Try alternative selectors for chapter content
	scraper.collector.OnHTML(".chapter-content", func(e *colly.HTMLElement) {
		logger.Infof("Found .chapter-content element")
		e.ForEach("p", func(_ int, paragraph *colly.HTMLElement) {
			chapterData.Paragraph = append(chapterData.Paragraph, paragraph.Text)
		})
		logger.Infof("Found %d paragraphs in chapter (alt selector)", len(chapterData.Paragraph))
	})

	scraper.collector.OnHTML(".content", func(e *colly.HTMLElement) {
		logger.Infof("Found .content element")
		e.ForEach("p", func(_ int, paragraph *colly.HTMLElement) {
			chapterData.Paragraph = append(chapterData.Paragraph, paragraph.Text)
		})
		logger.Infof("Found %d paragraphs in chapter (content selector)", len(chapterData.Paragraph))
	})

	// Try any link that might be the next chapter
	scraper.collector.OnHTML("a[href*='/chapter/']", func(e *colly.HTMLElement) {
		logger.Infof("Found chapter link: %s -> %s", e.Text, e.Attr("href"))
		if nextURL == "" && strings.Contains(e.Text, "Next") {
			nextURL = e.Attr("href")
			logger.Infof("Set next URL from link: %s", nextURL)
		}
	})

	defer func() {
		scraper.collector.OnHTMLDetach("#chr-content")
		scraper.collector.OnHTMLDetach("#next_chap")
		scraper.collector.OnHTMLDetach(".chapter-content")
		scraper.collector.OnHTMLDetach(".content")
		scraper.collector.OnHTMLDetach("body")
	}()

	scraper.collector.OnError(func(r *colly.Response, err error) {
		logger.Errorf("Request failed for URL %s: %v", r.Request.URL, err)
	})

	err := scraper.collector.Visit(url)
	if err != nil {
		logger.Errorf("Failed to visit URL %s: %v", url, err)
		return "", logger.Errorf("failed to visit url %s : %s", url, err.Error())
	}

	logger.Infof("Chapter scraping completed. Paragraphs: %d, Next URL: %s", len(chapterData.Paragraph), nextURL)

	if nextURL == "" {
		return "", logger.Errorf("Failed to get next url, current url is : %s", url)
	}

	return nextURL, nil
}

func NewBinNovelScrapper(app *core.App) NovelBinScraper {
	logger.Infof("Create %s's scraper", NovelBinScraperName)

	// Create collector with proper user agent and settings
	collector := colly.NewCollector(
		colly.AllowURLRevisit(),
		colly.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36"),
		colly.Async(true),
	)

	// Set random delays to mimic human behavior
	collector.Limit(&colly.LimitRule{
		DomainGlob:  "*",
		RandomDelay: 3 * time.Second,
		Parallelism: 1,
	})

	// Set up proper headers to avoid 403 Forbidden errors
	collector.OnRequest(func(r *colly.Request) {
		logger.Infof("Making request to: %s", r.URL)

		// Add random delay to mimic human behavior
		time.Sleep(time.Duration(1+rand.Intn(3)) * time.Second)

		r.Headers.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
		r.Headers.Set("Accept-Language", "en-US,en;q=0.9")
		r.Headers.Set("Accept-Encoding", "gzip, deflate, br")
		r.Headers.Set("Connection", "keep-alive")
		r.Headers.Set("Upgrade-Insecure-Requests", "1")
		r.Headers.Set("Cache-Control", "no-cache")
		r.Headers.Set("Pragma", "no-cache")
		r.Headers.Set("Sec-Fetch-Dest", "document")
		r.Headers.Set("Sec-Fetch-Mode", "navigate")
		r.Headers.Set("Sec-Fetch-Site", "none")
		r.Headers.Set("Sec-Fetch-User", "?1")
		r.Headers.Set("DNT", "1")
		r.Headers.Set("Sec-Ch-Ua", `"Not_A Brand";v="8", "Chromium";v="120", "Google Chrome";v="120"`)
		r.Headers.Set("Sec-Ch-Ua-Mobile", "?0")
		r.Headers.Set("Sec-Ch-Ua-Platform", `"Windows"`)

		logger.Infof("Making request to: %s", r.URL)
	})

	scraper := NovelBinScraper{
		collector: collector,
		app:       app,
	}

	return scraper
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
			logger.Infof("Attempting to scrape chapter %d (attempt %d/5): %s", i, retry+1, url)
			nextURL, err := scraper.scrapPage(url, &chapterData)
			if err == nil {
				logger.Infof("Successfully scraped chapter %d, next URL: %s", i, nextURL)
				url = nextURL
				break
			}
			logger.Errorf("Failed to get novel at URL %s for chapter %d (attempt %d/5): %v", url, i, retry+1, err)
			if retry < 4 { // Don't sleep on the last attempt
				sleepTime := time.Second * time.Duration(retry+1)
				logger.Infof("Retrying in %v...", sleepTime)
				time.Sleep(sleepTime)
			}
		}
		if retry == 5 {
			logger.Errorf("Failed to get novel at URL %s for chapter %d after 5 attempts", url, i)
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
