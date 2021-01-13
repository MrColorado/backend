package converter

import (
	"fmt"
	"os"

	"github.com/MrColorado/epubScraper/utils"
	"github.com/bmaupin/go-epub"
)

// EpubConverter convert novel data ton epub format
type EpubConverter struct{}

func (converter *EpubConverter) convertMetaData(e *epub.Epub, inputPath string, novelName string) {
	novelMetaData, err := utils.ImportMetaData(inputPath, novelName)
	if err != nil {
		println(err.Error())
		return
	}
	e.SetAuthor(novelMetaData.Author)
	summary := ""
	for _, paragraph := range novelMetaData.Summary {
		summary += fmt.Sprintf("<p>%s</p>", paragraph)
	}
	e.SetDescription(summary)

	coverImagePath, _ := e.AddImage(novelMetaData.ImagePath, "cover.png")
	e.SetCover(coverImagePath, "")
}

func (converter *EpubConverter) convertToNovel(inputPath string, outputPath string, novelName string, startChapter int, endChapter int) {
	e := epub.NewEpub(fmt.Sprintf("%s %d-%d", novelName, startChapter, endChapter))
	converter.convertMetaData(e, inputPath, novelName)

	for i := startChapter; i <= endChapter; i++ {
		novelChapterData, err := utils.ImportNovel(fmt.Sprintf("%s/%s/%04d.json", inputPath, novelName, i))
		if err != nil {
			println(err.Error())
			continue
		}

		bodySection := fmt.Sprintf("<h1>Chapter %d</h1>", novelChapterData.Chapter)
		for _, paragraph := range novelChapterData.Paragraph {
			bodySection += fmt.Sprintf("<p>%s</p>", paragraph)
		}
		if _, err := e.AddSection(bodySection, fmt.Sprintf("Chapter %d", i), "", ""); err != nil {
			fmt.Printf("Fail to add chapter %d of novel %s\n", i, novelName)
		}
	}

	directoryPath := fmt.Sprintf("%s/%s", outputPath, novelName)
	if _, err := os.Stat(directoryPath); os.IsNotExist(err) {
		if os.Mkdir(directoryPath, os.ModePerm) != nil {
			fmt.Printf("Failed to create directory : %s\n", directoryPath)
		}
	}

	e.Write(fmt.Sprintf("%s/%s/%s_%04d-%04d.epub", outputPath, novelName, novelName, startChapter, endChapter))
	fmt.Printf("%s/%s/%s_%04d-%04d.epub\n", outputPath, novelName, novelName, startChapter, endChapter)
}

// ConvertPartialNovel convert partial novel to epub (startChapter include / endChapter included)
func (converter *EpubConverter) ConvertPartialNovel(inputPath string, outputPath string, novelName string, startChapter int, endChapter int) {
	if endChapter > 100 && startChapter%100 != 1 {
		toModulo100 := 100 - startChapter%100
		converter.convertToNovel(inputPath, outputPath, novelName, startChapter, startChapter+toModulo100)
		startChapter += toModulo100
	}

	numberOfBook := (endChapter - startChapter) / 100
	firstBook := startChapter / 100

	for i := firstBook; i < firstBook+numberOfBook; i++ {
		converter.convertToNovel(inputPath, outputPath, novelName, i*100+1, (i+1)*100)
	}

	if endChapter%100 != 0 {
		converter.convertToNovel(inputPath, outputPath, novelName, (firstBook+numberOfBook)*100+1, endChapter)
	}
}

// ConvertNovel convert every novel in inputPath to epub format
func (converter *EpubConverter) ConvertNovel(inputPath string, outputPath string, novelName string) {
	numbreOfChapter := utils.NumberOfChapter(inputPath, novelName)
	fmt.Printf("Number of book : %d\n", numbreOfChapter)
	numberOfBook := numbreOfChapter / 100

	for i := 0; i < numberOfBook; i++ {
		converter.convertToNovel(inputPath, outputPath, novelName, i*100+1, (i+1)*100)
	}

	if numbreOfChapter%100 != 0 {
		converter.convertToNovel(inputPath, outputPath, novelName, numberOfBook*100+1, numbreOfChapter)
	}
}
