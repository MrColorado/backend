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

func (scraper *ReadNovelScraper) getNbChapter(c *colly.Collector, url string) int {
	nbChapter := 0

	c.OnHTML(".l-chapter", func(e *colly.HTMLElement) {
		splittedString := strings.Split(e.ChildAttr("a", "title"), " ")
		nbChapter, _ = strconv.Atoi(splittedString[1])
	})

	c.Visit(url)
	return nbChapter
}

func (scraper *ReadNovelScraper) scrapMetaData(c *colly.Collector, url string, novelName string) utils.NovelMetaData {
	fmt.Printf("Scrape : %s\n", url)
	novelMetaData := utils.NovelMetaData{
		Title: novelName,
		// ImagePath: fmt.Sprintf("%s/%s/cover", outputPath, novelName),
	}

	c.OnHTML(".btn-read-now", func(e *colly.HTMLElement) {
		novelMetaData.FirstChapterURL = readNovelURL + e.Attr("href")
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

func (scraper *ReadNovelScraper) scrapPage(c *colly.Collector, url string) (utils.NovelChapterData, string) {
	fmt.Printf("Scrape : %s\n", url)
	novelData := utils.NovelChapterData{}
	nextURL := ""

	c.OnHTML("#next_chap", func(e *colly.HTMLElement) {
		nextURL = readNovelURL + e.Attr("href")
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

	return novelData, nextURL
}

// ScrapeNovel get each chater of a specific novel
func (scraper *ReadNovelScraper) ScrapeNovel(c *colly.Collector, novelName string, outputPath string) {
	nbChapter := scraper.getNbChapter(colly.NewCollector(), fmt.Sprintf("%s/%s.html", readNovelURL, novelName))
	scraper.ScrapPartialNovel(c, novelName, outputPath, 1, nbChapter)
}

// ScrapPartialNovel get specified chapter of a novel
func (scraper *ReadNovelScraper) ScrapPartialNovel(c *colly.Collector, novelName string, outputPath string, startChapter int, endChapter int) {
	if utils.MataDataNotExist(fmt.Sprintf("%s/%s", outputPath, novelName)) {
		novelMetaData := scraper.scrapMetaData(c, fmt.Sprintf("%s/%s.html", readNovelURL, novelName), novelName)
		utils.ExportNovelMetaData(outputPath, novelName, novelMetaData)
	}

	novelMetaData, err := utils.ImportMetaData(fmt.Sprintf("%s/%s", outputPath, novelName))
	if err != nil {
		fmt.Println(err)
	}

	url := novelMetaData.FirstChapterURL
	if startChapter > 1 {
		url = novelMetaData.NextURL
	}

	for ; startChapter <= endChapter && strings.Compare(url, "") != 0; startChapter++ {
		var novel utils.NovelChapterData
		novel, url = scraper.scrapPage(c, url)
		novel.Chapter = startChapter
		novelMetaData.NextURL = url
		utils.ExportNovelChapter(outputPath, novelName, novel)
		utils.ExportNovelMetaData(outputPath, novelName, novelMetaData)
	}
}
