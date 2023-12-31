package converter

import (
	"bytes"
	"fmt"

	"github.com/MrColorado/epubScraper/configuration"
	"github.com/MrColorado/epubScraper/models"
	"github.com/MrColorado/epubScraper/utils"
	"github.com/bmaupin/go-epub"
)

// EpubConverter convert novel data ton epub format
type EpubConverter struct {
	io utils.S3IO
}

func NewEpubConverter(_ configuration.ConverterConfigStruct, io utils.S3IO) EpubConverter {
	return EpubConverter{
		io: io,
	}
}

func (converter EpubConverter) convertMetaData(e *epub.Epub, novelName string) error {
	data, err := converter.io.ImportMetaData(novelName)
	if err != nil {
		println(err.Error())
		return fmt.Errorf("failed to import metaData for novel %s", novelName)
	}

	e.SetAuthor(data.Author)
	summary := ""
	for _, paragraph := range data.Summary {
		summary += fmt.Sprintf("<p>%s</p>", paragraph)
	}
	e.SetDescription(summary)

	return nil
}

func (converter EpubConverter) convertToNovel(novelName string, startChapter int, endChapter int) error {
	fileName := fmt.Sprintf("%s-%04d-%04d", novelName, startChapter, endChapter)
	e := epub.NewEpub(fileName)
	converter.convertMetaData(e, novelName)

	for i := startChapter; i <= endChapter; i++ {
		chapterData, err := converter.io.ImportNovelChapter(novelName, i)
		chapterData.Chapter = 1
		if err != nil {
			fmt.Println(err.Error())
			continue
		}

		bodySection := fmt.Sprintf("<h1>Chapter %d</h1>", chapterData.Chapter)
		for _, paragraph := range chapterData.Paragraph {
			bodySection += fmt.Sprintf("<p>%s</p>", paragraph)
		}
		if _, err := e.AddSection(bodySection, fmt.Sprintf("Chapter %d", i), "", ""); err != nil {
			fmt.Printf("fail to add chapter %d of novel %s\n", i, novelName)
			continue
		}
	}

	buf := new(bytes.Buffer)
	_, err := e.WriteTo(buf)
	if err != nil {
		fmt.Println(err.Error())
		return fmt.Errorf("failed to write epub file %s in a buffer", fileName)
	}

	book := models.BookData{
		Start: startChapter,
		End:   endChapter,
	}
	err = converter.io.ExportBook(novelName, fileName, buf.Bytes(), book)
	if err != nil {
		fmt.Println(err.Error())
		return fmt.Errorf("failed to export epub file %s", fileName)
	}
	return nil
}

// ConvertPartialNovel convert partial novel to epub (startChapter include / endChapter included)
func (converter EpubConverter) ConvertPartialNovel(novelName string, startChapter int, endChapter int) error {
	if endChapter > 100 && startChapter%100 != 1 {
		toModulo100 := 100 - startChapter%100
		err := converter.convertToNovel(novelName, startChapter, startChapter+toModulo100)
		if err != nil {
			fmt.Println(err.Error())
		}
		startChapter += toModulo100
	}

	numberOfBook := (endChapter - startChapter) / 100
	firstBook := startChapter / 100

	for i := firstBook; i < firstBook+numberOfBook; i++ {
		err := converter.convertToNovel(novelName, i*100+1, (i+1)*100)
		if err != nil {
			fmt.Println(err.Error())
		}
	}

	if endChapter%100 != 0 {
		err := converter.convertToNovel(novelName, (firstBook+numberOfBook)*100+1, endChapter)
		if err != nil {
			fmt.Println(err.Error())
		}
	}
	return nil
}

// ConvertNovel convert every novel in inputPath to epub format
func (converter EpubConverter) ConvertNovel(novelName string) error {
	nbChapter, err := converter.io.NumberOfChapter(novelName)
	if err != nil {
		fmt.Println(err.Error())
		return fmt.Errorf("failed to get number of chapter for novel %s", novelName)
	}

	rest := 0
	if nbChapter%100 != 0 {
		rest += 1
	}
	nbBook := nbChapter / 100
	fmt.Printf("for novel %s there are %d chapter and so %d books\n", novelName, nbChapter, nbBook+rest)

	for i := 0; i < nbBook; i++ {
		err := converter.convertToNovel(novelName, (i*100)+1, (i+1)*100)
		if err != nil {
			fmt.Println(err.Error())
		}
	}

	if rest != 0 {
		err := converter.convertToNovel(novelName, (nbBook*100)+1, nbChapter)
		if err != nil {
			fmt.Println(err.Error())
		}
	}
	return nil
}
