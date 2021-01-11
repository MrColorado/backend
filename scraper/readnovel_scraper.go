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

func (scraper *ReadNovelScraper) scrapMetaData(c *colly.Collector, url string, outputPath string, novelName string) utils.NovelMetaData {
	novelMetaData := utils.NovelMetaData{
		Title:     novelName,
		ImagePath: fmt.Sprintf("%s/%s/cover", outputPath, novelName),
	}

	c.OnHTML(".l-chapter", func(e *colly.HTMLElement) {
		splittedString := strings.Split(e.ChildAttr("a", "title"), " ")
		novelMetaData.NumberOfChapter, _ = strconv.Atoi(splittedString[1])
	})

	c.OnHTML("#tab-description", func(e *colly.HTMLElement) {
		e.ForEach("p", func(_ int, paragraph *colly.HTMLElement) {
			novelMetaData.Summary = append(novelMetaData.Summary, paragraph.Text)
		})
	})

	c.OnHTML(".info-meta", func(e *colly.HTMLElement) {
		e.ForEach("li", func(_ int, li *colly.HTMLElement) {
			if title := li.ChildText("h3"); title == "Author:" {
				li.ForEach("a", func(_ int, a *colly.HTMLElement) {
					novelMetaData.Author = a.Text
					return
				})
			}
		})
	})

	c.Visit(url)
	return novelMetaData
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
func (scraper *ReadNovelScraper) ScrapPartialNovel(c *colly.Collector, novelName string, startChapter int, endChapter int, outputPath string) {
	if utils.MataDateAreExisting(fmt.Sprintf("%s/%s.html", readNovelURL, novelName)) {
		novelMetaData := scraper.scrapMetaData(c, fmt.Sprintf("%s/%s.html", readNovelURL, novelName), outputPath, novelName)
		utils.ExportNovelMetaData(outputPath, novelMetaData)
	}
	for ; startChapter <= endChapter; startChapter++ {
		novel, isNextPage := scraper.scrapPage(c, fmt.Sprintf("%s/%s/chapter-%d.html", readNovelURL, novelName, startChapter))
		novel.Chapter = startChapter
		utils.ExportNovelChapter(outputPath, novelName, novel)

		if isNextPage == false {
			break
		}
	}
}

// ScrapeNovel get each chater of a specific novel
func (scraper *ReadNovelScraper) ScrapeNovel(c *colly.Collector, novelName string, outputPath string) {
	nbChapter := 0
	if utils.MataDateAreExisting(fmt.Sprintf("%s/%s.html", readNovelURL, novelName)) {
		fmt.Println("FUCK")
		novelData, err := utils.ImportMetaData(fmt.Sprintf("%s/%s.html", readNovelURL, novelName))
		if err != nil {
			return
		}
		nbChapter = novelData.NumberOfChapter
	} else {
		fmt.Println(fmt.Sprintf("%s/%s.html", readNovelURL, novelName))
		novelMetaData := scraper.scrapMetaData(c, fmt.Sprintf("%s/%s.html", readNovelURL, novelName), outputPath, novelName)
		utils.ExportNovelMetaData(outputPath, novelMetaData)
		nbChapter = novelMetaData.NumberOfChapter
	}
	scraper.ScrapPartialNovel(c, novelName, 1, nbChapter, outputPath)
}
