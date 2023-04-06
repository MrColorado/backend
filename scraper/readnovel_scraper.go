package scraper

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/MrColorado/epubScraper/configuration"
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

func (scraper ReadNovelScraper) scrapMetaData(url string, novelMetaData *utils.NovelMetaData) {
	fmt.Printf("Scrape metaData : %s\n", url)

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

func (scraper ReadNovelScraper) scrapPage(url string, chapterData *utils.NovelChapterData) string {
	fmt.Printf("Scrape : %s\n", url)
	nextURL := ""

	scraper.collector.OnHTML("#next_chap", func(e *colly.HTMLElement) {
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

// ScrapeNovel get each chapter of a specific novel
func (scraper ReadNovelScraper) ScrapeNovel(novelName string) {
	scraper.ScrapeNovelStart(novelName, 1)
}

// ScrapeNovelStart get chapter of a specific novel starting a defined chapter
func (scraper ReadNovelScraper) ScrapeNovelStart(novelName string, startChapter int) {
	var metaData utils.NovelMetaData
	err := scraper.io.ImportMetaData(novelName, &metaData)
	if err != nil || metaData.Title == "" {
		scraper.scrapMetaData(fmt.Sprintf("%s/%s.html", readNovelURL, novelName), &metaData)
		scraper.io.ExportMetaData(novelName, metaData)
	}
	scraper.ScrapeNovelStartEnd(novelName, startChapter, metaData.NbChapter)
}

// ScrapeNovelStartEnd get chapter of specied novel between range
func (scraper ReadNovelScraper) ScrapeNovelStartEnd(novelName string, startChapter int, endChapter int) {
	if endChapter == -1 {
		endChapter = math.MaxInt
	}

	var metaData utils.NovelMetaData
	err := scraper.io.ImportMetaData(novelName, &metaData)
	if err != nil || metaData.Title == "" {
		scraper.scrapMetaData(fmt.Sprintf("%s/%s.html", readNovelURL, novelName), &metaData)
		scraper.io.ExportMetaData(novelName, metaData)
	}

	url := metaData.FirstChapterURL
	for i := 1; i <= endChapter && strings.Compare(url, "") != 0; i++ {
		chapterData := utils.NovelChapterData{
			Chapter: i,
		}
		url = scraper.scrapPage(url, &chapterData)
		metaData.NextURL = url

		if i >= startChapter {
			scraper.io.ExportNovelChapter(novelName, chapterData)
			scraper.io.ExportMetaData(novelName, metaData)
		}
	}
}
