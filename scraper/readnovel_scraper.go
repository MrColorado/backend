package scraper

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/MrColorado/epubScraper/configuration"
	"github.com/MrColorado/epubScraper/utils"

	"github.com/gocolly/colly"
)

const (
	readNovelURL string = "https://readnovelfull.com"
)

// ReadNovelScraper scrapper use to scrap https://vipnovel.com/ website
type ReadNovelScraper struct {
	collector *colly.Collector
	io        utils.IO
}

func (scraper ReadNovelScraper) getNbChapter(c *colly.Collector, url string) int {
	nbChapter := 0

	c.OnHTML(".l-chapter", func(e *colly.HTMLElement) {
		splittedString := strings.Split(e.ChildAttr("a", "title"), " ")
		nbChapter, _ = strconv.Atoi(splittedString[1])
	})

	c.Visit(url)
	return nbChapter
}

func (scraper ReadNovelScraper) scrapMetaData(url string, novelMetaData *utils.NovelMetaData) {
	fmt.Printf("Scrape metaData : %s\n", url)

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
					return
				})
			}
		})
	})
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

	scraper.collector.Visit(url)

	return nextURL
}

func NewReadNovelScrapper(_ configuration.ScraperConfigStruct, io utils.IO) ReadNovelScraper {
	return ReadNovelScraper{
		collector: colly.NewCollector(),
		io:        io,
	}
}

// ScrapeNovel get each chater of a specific novel
func (scraper ReadNovelScraper) ScrapeNovel(novelName string) {
	nbChapter := scraper.getNbChapter(colly.NewCollector(), fmt.Sprintf("%s/%s.html", readNovelURL, novelName))
	fmt.Printf("Number of chapter : %d\n", nbChapter)
	scraper.ScrapPartialNovel(novelName, 1, nbChapter)
}

// ScrapPartialNovel get specified chapter of a novel
func (scraper ReadNovelScraper) ScrapPartialNovel(novelName string, startChapter int, endChapter int) {
	var metaData utils.NovelMetaData
	err := scraper.io.ImportMetaData(novelName, &metaData)
	if err != nil || metaData.Title == "" {
		scraper.scrapMetaData(fmt.Sprintf("%s/%s.html", readNovelURL, novelName), &metaData)
		scraper.io.ExportMetaData(novelName, metaData)
	}

	url := metaData.NextURL
	if url == "" {
		url = metaData.FirstChapterURL
	}

	for ; startChapter <= endChapter && strings.Compare(url, "") != 0; startChapter++ {
		chapterData := utils.NovelChapterData{
			Chapter: startChapter,
		}
		url = scraper.scrapPage(url, &chapterData)
		metaData.NextURL = url

		scraper.io.ExportNovelChapter(novelName, chapterData)
		scraper.io.ExportMetaData(novelName, metaData)
	}
}
