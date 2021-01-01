package scraper

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/MrColorado/epubScraper/utils"

	"github.com/gocolly/colly"
)

const (
	readNovelURL string = "https://readnovelfull.com"
)

// ReadNovelScraper scrapper use to scrap https://vipnovel.com/ website
type ReadNovelScraper struct{}

func (scraper *ReadNovelScraper) scrapMainPage(c *colly.Collector, url string) int {
	endChapter := 0

	c.OnHTML(".l-chapter", func(e *colly.HTMLElement) {
		splittedString := strings.Split(e.ChildAttr("a", "title"), " ")
		endChapter, _ = strconv.Atoi(splittedString[1])
	})

	c.Visit(url)

	return endChapter
}

// ScrapeNovel get each chater of a specific novel
func (scraper *ReadNovelScraper) ScrapeNovel(c *colly.Collector, novelName string) {
	endChapter := scraper.scrapMainPage(c, fmt.Sprintf("%s/%s.html", readNovelURL, novelName))
	scraper.ScrapPartialNovel(c, novelName, 1, endChapter)
}

func (scraper *ReadNovelScraper) scrapPage(c *colly.Collector, url string) (utils.NovelChapterData, bool) {
	novelData := utils.NovelChapterData{}
	isNextPage := true

	c.OnHTML("#next_chap", func(e *colly.HTMLElement) {
		if e.Attr("disabled") != "" {
			isNextPage = false
		}
	})

	c.OnHTML("#chr-content", func(e *colly.HTMLElement) {
		e.ForEach("p", func(_ int, paragraph *colly.HTMLElement) {
			if title := paragraph.ChildText("span"); title != "" {
				novelData.Title = title
			} else {
				novelData.Paragraph = append(novelData.Paragraph, paragraph.Text)
			}
		})
	})

	c.Visit(url)

	return novelData, isNextPage
}

// ScrapPartialNovel get specified chapter of a novel
func (scraper *ReadNovelScraper) ScrapPartialNovel(c *colly.Collector, novelName string, startChapter int, endChapter int) {
	for ; startChapter <= endChapter; startChapter++ {
		novel, isNextPage := scraper.scrapPage(c, fmt.Sprintf("%s/%s/chapter-%04d.html", readNovelURL, novelName, startChapter))
		novel.Chapter = startChapter
		utils.ExportNovelChapter("/home/mrcolorado/Novels", novelName, novel)

		if isNextPage == false {
			break
		}
	}
}
