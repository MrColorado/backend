package scraper

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/MrColorado/epubScraper/configuration"
	"github.com/MrColorado/epubScraper/models"
	"github.com/MrColorado/epubScraper/utils"

	"github.com/gocolly/colly/v2"
)

const (
	readNovelURL string = "https://readnovelfull.com"
)

type ReadNovelScraper struct {
	collector *colly.Collector
	io        utils.IO
}

func (scraper ReadNovelScraper) findNovelUrl(novelName string) string {
	url := fmt.Sprintf("%s/novel-list/search?keyword=%s", readNovelURL, strings.ReplaceAll(novelName, " ", "+"))
	nbFound := 0
	novelUrl := ""

	scraper.collector.OnHTML(".list-novel", func(e *colly.HTMLElement) {
		e.ForEach(".novel-title", func(_ int, title *colly.HTMLElement) {
			nbFound += 1
			novelUrl = fmt.Sprintf("%s%s", readNovelURL, title.ChildAttr("a", "href"))
		})
	})

	scraper.collector.Visit(url)
	if nbFound == 1 {
		return novelUrl
	}
	return ""
}

func (scraper ReadNovelScraper) scrapMetaData(url string, novelMetaData *models.NovelMetaData) {
	fmt.Printf("Scrape metaData : %s\n", url)
	novelMetaData.CurrentChapter = 1

	scraper.collector.OnHTML(".l-chapter", func(e *colly.HTMLElement) {
		splittedString := strings.Split(e.ChildAttr("a", "title"), " ")
		nbChapter, err := strconv.Atoi(splittedString[1])
		if err != nil {
			novelMetaData.NbChapter = -1
		} else {
			novelMetaData.NbChapter = nbChapter
		}
	})

	scraper.collector.OnHTML(".title", func(e *colly.HTMLElement) {
		novelMetaData.Title = e.Text
	})

	scraper.collector.OnHTML(".btn-read-now", func(e *colly.HTMLElement) {
		novelMetaData.FirstChapterURL = readNovelURL + e.Attr("href")
	})

	scraper.collector.OnHTML("#tab-description", func(e *colly.HTMLElement) {
		e.ForEach("p", func(_ int, paragraph *colly.HTMLElement) {
			novelMetaData.Summary = append(novelMetaData.Summary, paragraph.Text)
		})
	})

	scraper.collector.OnHTML(".info-meta", func(e *colly.HTMLElement) {
		e.ForEach("li", func(_ int, li *colly.HTMLElement) {
			if title := li.ChildText("h3"); title == "Author:" {
				li.ForEach("a", func(_ int, a *colly.HTMLElement) {
					novelMetaData.Author = a.Text
				})
			}
		})
	})

	defer func() {
		scraper.collector.OnHTMLDetach("#tab-description")
		scraper.collector.OnHTMLDetach(".btn-read-now")
		scraper.collector.OnHTMLDetach(".l-chapter")
		scraper.collector.OnHTMLDetach(".next_chap")
		scraper.collector.OnHTMLDetach(".title")
	}()

	scraper.collector.Visit(url)
}

func (scraper ReadNovelScraper) scrapPage(url string, chapterData *models.NovelChapterData) string {
	fmt.Printf("Scrape : %s\n", url)
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

func NewReadNovelScrapper(_ configuration.ScraperConfigStruct, io utils.IO) ReadNovelScraper {
	return ReadNovelScraper{
		collector: colly.NewCollector(),
		io:        io,
	}
}

// ScrapeNovelStart get chapter of a specific novel starting a defined chapter
func (scraper ReadNovelScraper) scrapeNovelStart(novelName string, startChapter int) {
	data, _ := scraper.io.ImportMetaData(novelName)

	if data.ID == 0 {
		novelUrl := scraper.findNovelUrl(novelName)
		scraper.scrapMetaData(novelUrl, &data)
		if data.Title == "" {
			fmt.Printf("Failed to get page of novel %s\n", novelName)
			return
		}
		scraper.io.ExportMetaData(data.Title, &data)
	}

	i := 1
	url := data.FirstChapterURL
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

		scraper.io.ExportNovelChapter(novelName, chapterData)
		scraper.io.ExportMetaData(novelName, &data)
	}
}

// ScrapeNovel get each chapter of a specific novel
func (scraper ReadNovelScraper) ScrapeNovel(novelName string) {
	data, _ := scraper.io.ImportMetaData(novelName)
	if data.ID == 0 {
		scraper.scrapeNovelStart(novelName, 1)
	} else {
		scraper.scrapeNovelStart(novelName, data.CurrentChapter)
	}
}
