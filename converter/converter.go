package converter

import (
	"fmt"
	"os"
	"strings"

	"github.com/MrColorado/epubScraper/utils"

	epub "github.com/bmaupin/go-epub"
)

// Converter interface that convert novelData to e-tablet format
type Converter interface {
	ConvertNovel(inputPath string, outputPath string)
	ConvertPartialNovel(inputPath string, outputPath string, startChapter int, endChapter int)
}

// EpubConverter convert novel data ton epub format
type EpubConverter struct{}

func (converter *EpubConverter) convertToNovel(inputPath string, outputPath string, startChapter int, endChapter int) {
	novelPathWords := strings.Split(inputPath, "/")
	novelName := novelPathWords[len(novelPathWords)-1]
	e := epub.NewEpub(fmt.Sprintf("%s %d-%d", novelName, startChapter, endChapter))

	for i := startChapter; i <= endChapter; i++ {
		novelChapterData, err := utils.InportNovel(fmt.Sprintf("%s/%04d.json", inputPath, i))
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

	directoryPath := fmt.Sprintf("%s", outputPath)
	if _, err := os.Stat(directoryPath); os.IsNotExist(err) {
		if os.Mkdir(directoryPath, os.ModePerm) != nil {
			fmt.Printf("Failed to create directory : %s\n", directoryPath)
		}
	}

	fmt.Printf("Write %s/%s_%04d-%04d.epub", outputPath, novelPathWords[len(novelPathWords)-1], startChapter, endChapter)
	e.Write(fmt.Sprintf("%s/%s_%04d-%04d.epub", outputPath, novelPathWords[len(novelPathWords)-1], startChapter, endChapter))
}

// ConvertPartialNovel convert partial novel to epub (startChapter include / endChapter included)
func (converter *EpubConverter) ConvertPartialNovel(inputPath string, outputPath string, startChapter int, endChapter int) {
	if startChapter%100 != 0 {
		toModulo100 := 100 - startChapter%100
		converter.convertToNovel(inputPath, outputPath, startChapter, startChapter+toModulo100)
		startChapter += toModulo100
	}

	numberOfBook := (endChapter - startChapter) / 100
	firstBook := startChapter / 100

	for i := firstBook; i < firstBook+numberOfBook; i++ {
		converter.convertToNovel(inputPath, outputPath, i*100, (i+1)*100)
	}

	if numberOfBook >= 1 && endChapter%100 != 0 {
		converter.convertToNovel(inputPath, outputPath, (firstBook+numberOfBook)*100, endChapter)
	}
}

// ConvertNovel convert every novel in inputPath to epub format
func (converter *EpubConverter) ConvertNovel(inputPath string, outputPath string) {
	numbreOfChapter := utils.NumberOfChapter(inputPath)
	numberOfBook := numbreOfChapter / 100

	for i := 0; i < numberOfBook; i++ {
		converter.convertToNovel(inputPath, outputPath, i*100, (i+1)*100)
	}

	if numberOfBook >= 1 && numbreOfChapter%100 != 0 {
		converter.convertToNovel(inputPath, outputPath, numberOfBook*100, numbreOfChapter)
	}
}
